package scraper

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/bots-house/google-play-parser/internal/parser"
	"github.com/bots-house/google-play-parser/internal/ramda"
	"github.com/bots-house/google-play-parser/internal/shared"
	"github.com/bots-house/google-play-parser/models"
	sh "github.com/bots-house/google-play-parser/shared"
)

func Developer(ctx context.Context, client sh.HTTPClient, spec models.DeveloperSpec) ([]models.App, error) {
	if err := spec.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	path, isNumber := getDeveloperPath(spec.DevID)

	body, _, err := request(ctx, client, requestSpec{
		url: path,
		params: url.Values{
			"id": []string{spec.DevID},
			"hl": []string{spec.Lang},
			"gl": []string{spec.Country},
		},
	})
	if err != nil {
		return nil, err
	}

	parsed, err := parser.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("parse developer data: %w", err)
	}

	apps, token, err := parseDeveloperApps(parsed, isNumber)
	if err != nil {
		return nil, fmt.Errorf("developer apps not found: %w", err)
	}

	apps = processPages(ctx, client, pagesSpec{
		token: token,
		apps:  apps,
		count: spec.Count,
	})

	if !spec.Full {
		return apps, nil
	}

	return processFullDetail(ctx, client, apps...), nil
}

func getDeveloperPath(devID string) (string, bool) {
	path := getURL(storeURL)

	_, err := strconv.ParseFloat(devID, strconv.IntSize)
	if err != nil {
		return path + "/developer", false
	}

	return path + "/dev", true
}

func getDeveloperClusterMappings(isNumber bool) shared.ClusterMapping {
	if isNumber {
		return shared.ClusterMapping{
			Apps:  []any{"ds:3", 0, 1, 0, 21, 0},
			Token: []any{"ds:3", 0, 1, 0, 21, 1, 3, 1},
		}
	}

	return shared.ClusterMapping{
		Apps:  []any{"ds:3", 0, 1, 0, 22, 0},
		Token: []any{"ds:3", 0, 1, 0, 22, 1, 3, 1},
	}
}

func developerMappingPath(isNumber bool, paths ...any) []any {
	if isNumber {
		return paths[1:]
	}

	return paths
}

func parseDeveloperApps(parsed *shared.ParsedObject, isNumber bool) ([]models.App, string, error) {
	clusterMapping := getDeveloperClusterMappings(isNumber)

	mapping := &shared.AppMapping{
		Title: developerMappingPath(isNumber, 0, 3),
		AppID: developerMappingPath(isNumber, 0, 0, 0),
		URL: shared.MappingWithFunc[string, string]{
			Path: developerMappingPath(isNumber, 0, 10, 4, 2),
			Fun:  produceURL,
		},
		Icon:      developerMappingPath(isNumber, 0, 1, 3, 2),
		Developer: developerMappingPath(isNumber, 0, 14),
		Currency:  developerMappingPath(isNumber, 0, 8, 1, 0, 1),
		Price: shared.MappingWithFunc[float64, float64]{
			Path: developerMappingPath(isNumber, 0, 8, 1, 0, 0),
			Fun:  func(f float64) float64 { return f / 1000000 },
		},
		Free: shared.MappingWithFunc[float64, bool]{
			Path: developerMappingPath(isNumber, 0, 8, 1, 0, 0),
			Fun:  func(f float64) bool { return f == 0 },
		},
		Summary:   developerMappingPath(isNumber, 0, 13, 1),
		ScoreText: developerMappingPath(isNumber, 0, 4, 0),
		Score:     developerMappingPath(isNumber, 0, 4, 1),
	}

	rawApps, ok := ramda.Path(clusterMapping.Apps, parsed.Data).([]any)
	if !ok {
		return nil, "", fmt.Errorf("apps not found")
	}

	apps := produceRawApps(rawApps, mapping)
	if len(apps) == 0 {
		return nil, "", fmt.Errorf("parse app failed")
	}

	token, _ := ramda.Path(clusterMapping.Token, parsed.Data).(string)

	return apps, token, nil
}
