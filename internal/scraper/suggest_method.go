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

func Suggest(ctx context.Context, client sh.HTTPClient, spec models.SearchSpec) ([]string, error) {
	if err := spec.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	reqBody := url.Values{
		"f.req": []string{fmt.Sprintf(suggestBody, spec.Query)},
	}

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

	raw, ok := ramda.Path([]any{0, 0}, data).([]any)
	if !ok {
		return nil, fmt.Errorf("parse response body")
	}

	return shared.MapCheck(raw, func(entry any) (string, bool) {
		data, ok := ramda.Path([]any{0}, entry).(string)
		return data, ok
	}), nil
}
