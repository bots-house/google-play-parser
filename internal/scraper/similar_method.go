package scraper

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/bots-house/google-play-parser/internal/parser"
	"github.com/bots-house/google-play-parser/internal/ramda"
	"github.com/bots-house/google-play-parser/internal/shared"
	"github.com/bots-house/google-play-parser/models"
	sh "github.com/bots-house/google-play-parser/shared"
)

func Similar(ctx context.Context, client sh.HTTPClient, spec models.ApplicationSpec) ([]models.App, error) {
	if err := spec.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	body, _, err := request(ctx, client, requestSpec{
		url: getURL(appsDetailsURL),
		params: url.Values{
			"id": []string{spec.AppID},
			"hl": []string{spec.Lang},
		},
	})
	if err != nil {
		return nil, err
	}

	parsed, err := parser.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	similarApps, err := parseSimilarApps(ctx, client, *parsed, spec)
	if err != nil {
		return nil, err
	}

	return similarApps, nil
}

func parseSimilarApps(ctx context.Context, client sh.HTTPClient, parsed shared.ParsedObject, spec models.ApplicationSpec) ([]models.App, error) {
	extracted := parser.ExtractDataWithServiceRequestID(parsed, clusterSpec)

	extractedClusters, ok := extracted.([]any)
	if !ok {
		return nil, fmt.Errorf("similar apps not found")
	}

	if len(extractedClusters) == 0 {
		return nil, fmt.Errorf("clusters not found")
	}

	clusterURL, ok := ramda.Path(clusterMapping.URL, extractedClusters[0]).(string)
	if !ok {
		return nil, fmt.Errorf("cluster url not found")
	}

	body, _, err := request(ctx, client, requestSpec{
		url: getURL(clusterURL),
	})
	if err != nil {
		return nil, err
	}

	similar, err := parser.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("no similar apps founded")
	}

	apps, err := processFirstPage(*similar, &clusterMapping)
	if err != nil {
		return nil, err
	}

	if !spec.Full {
		return apps, nil
	}

	return processFullDetail(ctx, client, apps...), nil
}

func processFirstPage(
	parsed shared.ParsedObject,
	mappings *shared.ClusterMapping,
) ([]models.App, error) {
	mapping := &shared.Mapping{
		Title: []any{3},
		AppID: []any{0, 0},
		URL: shared.MappingWithFunc[string, string]{
			Path: []any{10, 4, 2},
			Fun: func(u string) string {
				result, err := url.Parse(getURL(u))
				if err != nil {
					return ""
				}

				return result.String()
			},
		},
		Icon:      []any{1, 3, 2},
		Developer: []any{14},
		Currency:  []any{8, 1, 0, 1},
		Price: shared.MappingWithFunc[float64, float64]{
			Path: []any{8, 1, 0, 0},
			Fun:  func(price float64) float64 { return price / 1000000 },
		},
		Free: shared.MappingWithFunc[float64, bool]{
			Path: []any{8, 1, 0, 0},
			Fun:  func(f float64) bool { return f == 0 },
		},
		Summary:   []any{13, 1},
		ScoreText: []any{4, 0},
		Score:     []any{4, 1},
	}

	rawApps := ramda.Path(mappings.Apps, parsed.Data)

	rawAppsSlice, ok := rawApps.([]any)
	if !ok {
		return nil, fmt.Errorf("apps not found")
	}

	apps := produceRawApps(rawAppsSlice, mapping)
	if len(apps) == 0 {
		return nil, fmt.Errorf("no apps found")
	}

	return apps, nil
}

func produceRawApps(appsData []any, mapping *shared.Mapping) []models.App {
	apps := make([]models.App, 0)

	for _, appData := range appsData {
		app, ok := parser.Extract(appData, mapping)
		if !ok {
			log.Debug().Msg("app not found")
			continue
		}

		apps = append(apps, app)
	}

	return apps
}

func processFullDetail(ctx context.Context, client sh.HTTPClient, apps ...models.App) []models.App {
	var wg sync.WaitGroup

	wg.Add(len(apps))

	for idx := range apps {
		app := apps[idx]

		go func(idx int, app models.App) {
			defer wg.Done()
			fullApp, err := App(ctx, client, models.ApplicationSpec{
				AppID: app.AppID,
			})
			if err != nil {
				log.Error().Err(err).Msg("cannot produce full detail for app")
				return
			}

			apps[idx] = app.Assign(&fullApp)
		}(idx, app)

	}

	wg.Wait()

	return apps
}
