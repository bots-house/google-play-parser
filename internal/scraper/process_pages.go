package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/bots-house/google-play-parser/internal/ramda"
	"github.com/bots-house/google-play-parser/internal/shared"
	"github.com/bots-house/google-play-parser/models"
	sh "github.com/bots-house/google-play-parser/shared"
)

const (
	pagesURL  = `/_/PlayStoreUi/data/batchexecute?rpcids=qnKhOb&f.sid=-697906427155521722&bl=boq_playuiserver_20190903.08_p0&authuser&soc-app=121&soc-platform=1&soc-device=1&_reqid=1065213`
	pagesBody = `[[["qnKhOb","[[null,[[10,[10,%d]],true,null,[96,27,4,8,57,30,110,79,11,16,49,1,3,9,12,104,55,56,51,10,34,77]],null,\"%s\"]]",null,"generic"]]]`
)

type pagesSpec struct {
	token string
	count int
	apps  []models.App
}

func processPages(ctx context.Context, client sh.HTTPClient, spec pagesSpec) []models.App {
	if len(spec.apps) >= spec.count || spec.token == "" {
		return spec.apps[:spec.count]
	}

	reqBody := url.Values{
		"f.req": []string{fmt.Sprintf(pagesBody, spec.count, spec.token)},
	}

	body, _, err := request(ctx, client, requestSpec{
		url:    getURL(pagesURL),
		method: http.MethodPost,
		headers: http.Header{
			"Content-Type": []string{"application/x-www-form-urlencoded;charset=UTF-8"},
		},
		body: strings.NewReader(reqBody.Encode()),
	})
	if err != nil {
		log.Error().Err(err).Msg("process next pages")
		return spec.apps
	}

	body = body[6:]

	var data []any

	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	data, ok := data[0].([]any)
	if !ok {
		log.Error().Msg("apps data not found")
		return spec.apps
	}

	if len(data) < 3 {
		log.Error().Msg("apps data not found")
		return spec.apps
	}

	raw, ok := data[2].(string)
	if !ok {
		log.Error().Msg("apps data not found")
		return spec.apps
	}

	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		log.Error().Err(err).Msg("apps data not found")
		return spec.apps
	}

	// apps, err :=

	return processNextPages(ctx, client, data, spec)
}

func processNextPages(ctx context.Context, client sh.HTTPClient, rawData []any, spec pagesSpec) []models.App {
	rawApps, ok := ramda.Path([]any{0, 0, 0}, rawData).([]any)
	if !ok {
		log.Error().Msg("apps data not found")
		return spec.apps
	}

	mapping := &shared.AppMapping{
		Title: []any{2},
		AppID: []any{12, 0},
		URL: shared.MappingWithFunc[string, string]{
			Path: []any{9, 4, 2},
			Fun:  produceURL,
		},
		Icon:      []any{1, 1, 0, 3, 2},
		Developer: []any{4, 0, 0, 0},
		DeveloperID: shared.MappingWithFunc[string, string]{
			Path: []any{4, 0, 0, 1, 4, 2},
			Fun:  func(s string) string { return strings.Split(s, "id=")[1] },
		},
		PriceText: shared.MappingWithFunc[string, string]{
			Path: []any{7, 0, 3, 2, 1, 0, 2},
			Fun: func(s string) string {
				if s == "" {
					return "FREE"
				}

				return s
			},
		},
		Price: shared.MappingWithFunc[float64, float64]{
			Path: []any{7, 0, 3, 2, 1, 0, 1},
			Fun:  func(f float64) float64 { return f / 1000000 },
		},
		Summary:   []any{4, 1, 1, 1, 1},
		Score:     []any{6, 0, 2, 1, 1},
		ScoreText: []any{6, 0, 2, 1, 0},
	}

	apps := produceRawApps(rawApps, mapping)

	token, ok := ramda.Path([]any{0, 0, 7, 1}, rawData).(string)
	if !ok || token == "" {
		return append(spec.apps, apps...)
	}

	spec.token = token
	spec.apps = append(spec.apps, apps...)

	return processPages(ctx, client, spec)
}
