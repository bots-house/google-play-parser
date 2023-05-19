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

func Permissions(ctx context.Context, client sh.HTTPClient, spec models.ApplicationSpec) ([]models.Permission, error) {
	if err := spec.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	reqBody := url.Values{
		"f.req": []string{fmt.Sprintf(permissionsBody, spec.AppID)},
	}

	body, _, err := request(ctx, client, requestSpec{
		url:    getURL(permissionsURL),
		method: http.MethodPost,
		params: url.Values{
			"hl": []string{spec.Lang},
		},
		headers: http.Header{
			"Content-Type": []string{"application/x-www-form-urlencoded;charset=UTF-8"},
		},
		body: strings.NewReader(reqBody.Encode()),
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

	if spec.Full {
		return parseFullPermissions(data)
	}

	return parseShortPermissions(data)
}

func parseShortPermissions(data []any) ([]models.Permission, error) {
	commonPerms, ok := data[0].([]any)
	if !ok {
		return nil, fmt.Errorf("common perms not found")
	}

	commonPerms = shared.Filter(commonPerms, func(perm any) bool {
		slice, ok := perm.([]any)
		if !ok || len(slice) == 0 {
			return false
		}

		return true
	})

	permissionNames := shared.MapCheck(commonPerms, func(perm any) (string, bool) {
		perms, ok := perm.([]any)
		if !ok {
			return "", false
		}

		entry, ok := perms[0].(string)
		if !ok {
			return "", false
		}

		return entry, true
	})

	return shared.Map(permissionNames, func(perm string) models.Permission {
		return models.Permission{Type: perm}
	}), nil
}

func parseFullPermissions(data []any) ([]models.Permission, error) {
	filteredPerms := shared.Filter(data, func(entry any) bool {
		slice, ok := entry.([]any)
		if !ok || len(slice) == 0 {
			return false
		}

		return true
	})

	if len(filteredPerms) > 2 {
		filteredPerms = filteredPerms[:2]
	}

	return shared.Chain(filteredPerms, func(entry any) []models.Permission {
		slice, ok := entry.([]any)
		if !ok || len(slice) == 0 {
			return nil
		}

		return shared.Chain(slice, func(entry any) []models.Permission {
			slice, ok := entry.([]any)
			if !ok {
				return nil
			}

			typ, ok := slice[0].(string)
			if !ok {
				return nil
			}

			perms, ok := slice[2].([]any)
			if !ok {
				return nil
			}

			return shared.MapCheck(perms, func(entry any) (perm models.Permission, ok bool) {
				perm.Type = typ
				
				summary, ok := ramda.Path([]any{1}, entry).(string)
				if ok {
					perm.Summary = summary
				}

				return
			})
		})
	}), nil
}
