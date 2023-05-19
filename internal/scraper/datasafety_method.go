package scraper

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bots-house/google-play-parser/internal/parser"
	"github.com/bots-house/google-play-parser/internal/ramda"
	"github.com/bots-house/google-play-parser/internal/shared"
	"github.com/bots-house/google-play-parser/models"
	sh "github.com/bots-house/google-play-parser/shared"
)

func DataSafety(ctx context.Context, client sh.HTTPClient, spec models.ApplicationSpec) (models.DataSafety, error) {
	if err := spec.Validate(); err != nil {
		return models.DataSafety{}, fmt.Errorf("validation: %w", err)
	}

	body, _, err := request(ctx, client, requestSpec{
		url: getURL(datasafetyURL),
		params: url.Values{
			"id": []string{spec.AppID},
			"hl": []string{spec.Lang},
		},
	})
	if err != nil {
		return models.DataSafety{}, err
	}

	parsed, err := parser.Parse(body)
	if err != nil {
		return models.DataSafety{}, fmt.Errorf("data_safety parse: %w", err)
	}

	extracted, ok := parser.Extract[dataSafety](parsed.Data, shared.DataSafetyMapping{
		SharedData: shared.MappingWithFunc[any, []map[string]any]{
			Path: []any{"ds:3", 1, 2, 137, 4, 0, 0},
			Fun:  mapDataSafetyEntries,
		},
		CollectedData: shared.MappingWithFunc[any, []map[string]any]{
			Path: []any{"ds:3", 1, 2, 137, 4, 1, 0},
			Fun:  mapDataSafetyEntries,
		},
		PrivacyPolicyURL: []any{"ds:3", 1, 2, 99, 0, 5, 2},
		SecurityPractices: shared.MappingWithFunc[any, []map[string]any]{
			Path: []any{"ds:3", 1, 2, 137, 9, 2},
			Fun: func(a any) []map[string]any {
				slice, ok := a.([]any)
				if !ok || len(slice) == 0 {
					return nil
				}

				return shared.Map(slice, func(entry any) map[string]any {
					slice, ok := entry.([]any)
					if !ok || len(slice) == 0 {
						return nil
					}

					return map[string]any{
						"practice":    ramda.Path([]any{1}, entry),
						"description": ramda.Path([]any{2, 1}, entry),
					}
				})
			},
		},
	})
	if !ok {
		return models.DataSafety{}, fmt.Errorf("datasafety entry not found")
	}

	return extracted.toModel(), nil
}
