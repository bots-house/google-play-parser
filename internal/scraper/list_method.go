package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/bots-house/google-play-parser/internal/ramda"
	"github.com/bots-house/google-play-parser/internal/shared"
	"github.com/bots-house/google-play-parser/models"
	sh "github.com/bots-house/google-play-parser/shared"
)

func List(ctx context.Context, client sh.HTTPClient, spec models.ListSpec) ([]models.App, error) {
	if err := spec.Validate(); err != nil {
		return nil, err
	}

	params := url.Values{
		"hl": []string{spec.Lang},
		"gl": []string{spec.Country},
	}

	if spec.Age != "" {
		params.Set("age", spec.Age)
	}

	collection, ok := listClusterNames[spec.Collection]
	if !ok {
		return nil, fmt.Errorf("invalid collection")
	}

	requestBodyParams := fmt.Sprintf(listRequestBody, spec.Count, collection, spec.Category)

	requestBody := url.Values{
		"f.req": []string{requestBodyParams},
		"at":    []string{"AFSRYlx8XZfN8-O-IKASbNBDkB6T:1655531200971"},
	}

	body, _, err := request(ctx, client, requestSpec{
		method: http.MethodPost,
		url:    getURL(listURL),
		params: params,
		headers: http.Header{
			"Content-Type": []string{"application/x-www-form-urlencoded;charset=UTF-8"},
		},
		body: strings.NewReader(requestBody.Encode()),
	})
	if err != nil {
		return nil, err
	}

	bodyRaw := strings.Split(string(body), "\n")[2]

	var data []any

	if err := json.Unmarshal([]byte(bodyRaw), &data); err != nil {
		return nil, err
	}

	stringData, ok := ramda.Path([]any{0, 2}, data).(string)
	if !ok {
		return nil, fmt.Errorf("apps not found")
	}

	if err := json.Unmarshal([]byte(stringData), &data); err != nil {
		return nil, err
	}

	apps, err := parseCollectionApps(data)
	if err != nil {
		return nil, err
	}

	if !spec.Full {
		return apps, nil
	}

	return processFullDetail(ctx, client, apps...), nil
}

func parseCollectionApps(rawData []any) ([]models.App, error) {
	mapping := &shared.Mapping{
		Title: []any{0, 3},
		AppID: []any{0, 0, 0},
		URL: shared.MappingWithFunc[string, string]{
			Path: []any{0, 10, 4, 2},
			Fun: func(s string) string {
				u, err := url.Parse(getURL(s))
				if err != nil {
					return ""
				}

				return u.String()
			},
		},
		Icon:      []any{0, 1, 3, 2},
		Developer: []any{0, 14},
		Currency:  []any{0, 8, 1, 0, 1},
		Price: shared.MappingWithFunc[float64, float64]{
			Path: []any{0, 8, 1, 0, 0},
			Fun:  func(f float64) float64 { return f / 1000000 },
		},
		Free: shared.MappingWithFunc[float64, bool]{
			Path: []any{0, 8, 1, 0, 0},
			Fun:  func(f float64) bool { return f == 0 },
		},
		Summary:   []any{0, 13, 1},
		Score:     []any{0, 4, 1},
		ScoreText: []any{0, 4, 0},
	}

	rawApps, ok := ramda.Path([]any{0, 1, 0, 28, 0}, rawData).([]any)
	if !ok {
		return nil, fmt.Errorf("no apps data")
	}

	apps := produceRawApps(rawApps, mapping)
	if len(apps) == 0 {
		return nil, fmt.Errorf("parse apps failed")
	}

	return apps, nil
}
