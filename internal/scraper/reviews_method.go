package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/bots-house/google-play-parser/internal/parser"
	"github.com/bots-house/google-play-parser/internal/ramda"
	"github.com/bots-house/google-play-parser/internal/shared"
	"github.com/bots-house/google-play-parser/models"
	sh "github.com/bots-house/google-play-parser/shared"
)

func Reviews(ctx context.Context, client sh.HTTPClient, spec models.ReviewsSpec) ([]models.Review, error) {
	if err := spec.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	return produceReviewsRequest(ctx, client, nil, spec)
}

func produceReviewsRequest(ctx context.Context, client sh.HTTPClient, reviews []models.Review, spec models.ReviewsSpec) ([]models.Review, error) {
	body, _, err := request(ctx, client, requestSpec{
		url:    getURL(permissionsURL),
		method: http.MethodPost,
		headers: http.Header{
			"Content-Type": []string{"application/x-www-form-urlencoded;charset=UTF-8"},
		},
		params: url.Values{
			"hl": []string{spec.Lang},
			"gl": []string{spec.Country},
		},
		body: strings.NewReader(reviewsRequestBody(spec).Encode()),
	})
	if err != nil {
		return nil, err
	}

	body = body[6:]

	var data []any

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parse response body")
	}

	rawData, ok := ramda.Path([]any{0, 2}, data).(string)
	if !ok {
		return nil, fmt.Errorf("parse response body")
	}

	if err := json.Unmarshal([]byte(rawData), &data); err != nil {
		return nil, fmt.Errorf("parse response body")
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("unimplemented")
	}

	return parseReviewsResponse(ctx, client, data, reviews, spec)
}

func reviewsRequestBody(spec models.ReviewsSpec) url.Values {
	if spec.RequestType == "initial" {
		return url.Values{
			"f.req": []string{
				fmt.Sprintf(
					`[[["UsvDTd","[null,null,[2,%d,[%d,null,null],null,[]],[\"%s\",7]]",null,"generic"]]]`,
					spec.Count,
					spec.Count,
					spec.AppID,
				),
			},
		}
	}

	return url.Values{
		"f.req": []string{
			fmt.Sprintf(
				`[[["UsvDTd","[null,null,[2,%d,[%d,null,\"%s\"],null,[]],[\"%s\",7]]",null,"generic"]]]`,
				spec.Sort,
				spec.Count,
				spec.NextToken,
				spec.AppID,
			),
		},
	}
}

func parseReviewsResponse(
	ctx context.Context,
	client sh.HTTPClient,
	data []any,
	reviews []models.Review,
	spec models.ReviewsSpec,
) ([]models.Review, error) {
	if len(data) == 0 && spec.NextToken == "" && len(reviews) > 0 {
		return reviews[:spec.Count], nil
	}

	if len(data) == 0 && spec.NextToken == "" {
		return nil, fmt.Errorf("reviews data not found")
	}

	reviewsData, ok := data[0].([]any)
	if !ok {
		return nil, fmt.Errorf("reviews data not found")
	}

	mapping := &shared.ReviewsMapping{
		ID:        []any{0},
		UserName:  []any{1, 0},
		UserImage: []any{1, 1, 3, 2},
		Date: shared.MappingWithFunc[[]any, time.Time]{
			Path: []any{5},
			Fun:  reviewsDate,
		},
		Score: []any{2},
		ScoreText: shared.MappingWithFunc[float64, string]{
			Path: []any{2},
			Fun:  func(f float64) string { return strconv.FormatFloat(f, 'f', 2, strconv.IntSize) },
		},
		URL: shared.MappingWithFunc[string, string]{
			Path: []any{0},
			Fun: func(s string) string {
				path := fmt.Sprintf("/store/apps/details?id=%s&reviewId=%s", spec.AppID, s)
				u, err := url.Parse(getURL(path))
				if err != nil {
					return ""
				}

				return u.String()
			},
		},
		Summary: []any{4},
		ReplyDate: shared.MappingWithFunc[[]any, time.Time]{
			Path: []any{7, 2},
			Fun:  reviewsDate,
		},
		ReplyText: []any{7, 1},
		Version:   []any{10},
		Criteria: shared.MappingWithFunc[[]any, map[string]float64]{
			Path: []any{12, 0},
			Fun: func(a []any) map[string]float64 {
				result := make(map[string]float64)

				for _, entry := range a {
					key, ok := ramda.Path([]any{0}, entry).(string)
					if !ok {
						continue
					}

					value, _ := ramda.Path([]any{1, 0}, entry).(float64)

					result[key] = value
				}

				return result
			},
		},
	}

	reviews = append(reviews, parseRawReviews(reviewsData, mapping)...)

	token, ok := ramda.Path([]any{1, 1}, data).(string)
	if !ok || len(reviews) > spec.Count {
		return reviews[:spec.Count], nil
	}

	spec.NextToken = token

	return produceReviewsRequest(ctx, client, reviews, spec)
}

func parseRawReviews(rawReviews []any, mapping *shared.ReviewsMapping) []models.Review {
	result := make([]review, 0, len(rawReviews))

	for _, entry := range rawReviews {
		review, ok := parser.Extract[review](entry, mapping)
		if !ok {
			continue
		}

		result = append(result, review)
	}

	return shared.Map(result, func(review review) models.Review {
		return models.Review{
			ID:        uuid.MustParse(review.ID),
			URL:       review.URL,
			Title:     review.Title,
			Summary:   review.Summary,
			Score:     review.Score,
			ScoreText: review.ScoreText,
			UserName:  review.UserName,
			UserImage: review.UserImage,
			Version:   review.Version,
			Date:      review.Date,
			ReplyText: review.ReplyText,
			ReplyDate: review.ReplyDate,
			TumbsUp:   review.TumbsUp,
			Criteria: func(criteria map[string]float64) []models.ReviewCriteria {
				result := make([]models.ReviewCriteria, 0, len(criteria))

				for key, val := range criteria {
					result = append(result, models.ReviewCriteria{
						Criteria: key,
						Rating:   val,
					})
				}
				return result
			}(review.Criteria),
		}
	})
}
