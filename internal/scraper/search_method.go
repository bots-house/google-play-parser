package scraper

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/bots-house/google-play-parser/internal/parser"
	"github.com/bots-house/google-play-parser/internal/ramda"
	"github.com/bots-house/google-play-parser/internal/shared"
	"github.com/bots-house/google-play-parser/models"
	sh "github.com/bots-house/google-play-parser/shared"
)

func Search(ctx context.Context, client sh.HTTPClient, spec models.SearchSpec) ([]models.App, error) {
	if err := spec.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	validPrice := func(price string) string {
		switch strings.ToLower(price) {
		case "free":
			return "1"
		case "paid":
			return "2"
		default:
			return "0"
		}
	}

	body, _, err := request(ctx, client, requestSpec{
		url: getURL(searchURL),
		params: url.Values{
			"c":     []string{"apps"},
			"q":     []string{spec.Query},
			"hl":    []string{spec.Lang},
			"gl":    []string{spec.Country},
			"price": []string{validPrice(spec.Price)},
		},
	})
	if err != nil {
		return nil, err
	}

	parsed, err := parser.Parse(body)
	if err != nil {
		return nil, err
	}

	apps, err := processSearchResult(parsed, shared.ClusterMapping{
		Apps:     []any{"ds:4", 0, 1, 0, 23},
		Sections: []any{"ds:4", 0, 1},
	})
	if err != nil {
		return nil, err
	}

	if !spec.Full {
		return apps, nil
	}

	return processFullDetail(ctx, client, apps...), nil
}

func processSearchResult(parsed *shared.ParsedObject, initialMapping shared.ClusterMapping) ([]models.App, error) {
	mainAppMapping := &shared.AppMapping{
		Title: []any{16, 2, 0, 0},
		AppID: []any{16, 11, 0, 0},
		URL: shared.MappingWithFunc[string, string]{
			Path: []any{17, 0, 0, 4, 2},
			Fun:  produceURL,
		},
		Icon:      []any{16, 2, 95, 0, 3, 2},
		Developer: []any{16, 2, 68, 0},
		DeveloperID: shared.MappingWithFunc[string, string]{
			Path: []any{16, 2, 68, 1, 4, 2},
			Fun:  func(s string) string { return strings.Split(s, "?id=")[1] },
		},
		Currency: []any{17, 0, 2, 0, 1, 0, 1},
		Price: shared.MappingWithFunc[float64, float64]{
			Path: []any{17, 0, 2, 0, 1, 0, 0},
			Fun:  func(f float64) float64 { return f / 1000000 },
		},
		Free: shared.MappingWithFunc[float64, bool]{
			Path: []any{17, 0, 2, 0, 1, 0, 0},
			Fun:  func(f float64) bool { return f == 0 },
		},
		Summary:   []any{16, 2, 73, 0, 1},
		ScoreText: []any{16, 2, 51, 0, 0},
		Score:     []any{16, 2, 51, 0, 1},
	}

	moreResultMapping := &shared.AppMapping{
		Title: []any{0, 3},
		AppID: []any{0, 0, 0},
		URL: shared.MappingWithFunc[string, string]{
			Path: []any{0, 10, 4, 2},
			Fun:  produceURL,
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
		ScoreText: []any{0, 4, 0},
		Score:     []any{0, 4, 1},
	}

	sections, ok := ramda.Path(initialMapping.Sections, parsed.Data).([]any)
	if !ok {
		return nil, fmt.Errorf("sections not found")
	}

	moreSections := shared.Filter(sections, func(value any) bool {
		sectionTitle, ok := ramda.Path([]any{22, 1, 0}, value).(string)
		return ok && sectionTitle == ""
	})

	if len(moreSections) == 0 {
		return nil, fmt.Errorf("sections not found")
	}

	rawApps, ok := ramda.Path([]any{22, 0}, moreSections[0]).([]any)
	if !ok {
		return nil, fmt.Errorf("apps data not found")
	}

	apps := produceRawApps(rawApps, moreResultMapping)
	if len(apps) == 0 {
		return nil, fmt.Errorf("apps not found")
	}

	mainAppSection, ok := ramda.Path(initialMapping.Apps, parsed.Data).([]any)
	if !ok {
		return apps, nil
	}

	mainApp, ok := parser.Extract[models.App](mainAppSection, mainAppMapping)
	if !ok {
		log.Debug().Msg("main app data not found")
		return apps, nil
	}

	return append([]models.App{mainApp}, apps...), nil
}
