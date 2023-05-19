package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bots-house/google-play-parser/internal/ramda"
	"github.com/bots-house/google-play-parser/internal/shared"
	sh "github.com/bots-house/google-play-parser/shared"
)

func getURL(path string) string {
	return baseURL + path
}

func produceURL(path string) string {
	u, err := url.Parse(getURL(path))
	if err != nil {
		return ""
	}

	return u.String()
}

type requestSpec struct {
	method  string
	url     string
	params  url.Values
	headers http.Header
	body    io.Reader
}

func (spec *requestSpec) ensureNotNil() {
	if spec.method == "" {
		spec.method = http.MethodGet
	}
}

func (spec requestSpec) validate() error {
	if spec.url == "" {
		return fmt.Errorf("url required")
	}

	return nil
}

func request(ctx context.Context, client sh.HTTPClient, spec requestSpec) (body []byte, rawURL string, err error) {
	spec.ensureNotNil()

	if err := spec.validate(); err != nil {
		return nil, "", err
	}

	var requestBody io.Reader = http.NoBody
	if spec.body != nil {
		requestBody = spec.body
	}

	request, err := http.NewRequestWithContext(ctx, spec.method, spec.url, requestBody)
	if err != nil {
		return nil, "", fmt.Errorf("prepare request: %w", err)
	}

	if spec.headers != nil {
		request.Header = spec.headers
	}

	if spec.params != nil {
		request.URL.RawQuery = spec.params.Encode()
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, "", fmt.Errorf("do request: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected response: status: %s", response.Status)
	}

	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read body: %w", err)
	}

	return body, request.URL.String(), nil
}

func parseDescription(description []any) string {
	translation, _ := ramda.Path([]any{12, 0, 0, 1}, description).(string)
	original, _ := ramda.Path([]any{72, 0, 1}, description).(string)

	if translation == "" {
		return original
	}

	return translation
}

func mapDataSafetyEntries(raw any) []map[string]any {
	slice, ok := raw.([]any)
	if !ok || len(slice) == 0 {
		return nil
	}

	return shared.Map(slice, func(entry any) map[string]any {
		typ := ramda.Path([]any{1, 4}, entry)

		data, ok := ramda.Path([]any{4}, entry).([]any)
		if !ok {
			return nil
		}

		result := shared.Map(data, func(entry any) map[string]any {
			return map[string]any{
				"data":     ramda.Path([]any{0}, entry),
				"optional": ramda.Path([]any{1}, entry),
				"purpose":  ramda.Path([]any{2}, entry),
				"type":     typ,
			}
		})

		if len(result) == 0 {
			return nil
		}

		return result[0]
	})
}

func safeMapIndex[O any](m map[string]any, key string) (out O) {
	entry, ok := m[key]
	if !ok {
		return out
	}

	result, ok := entry.(O)
	if !ok {
		return out
	}

	return result
}

func reviewsDate(dateArray []any) time.Time {
	mills, ok := ramda.Path([]any{1}, dateArray).(float64)
	if !ok {
		mills = 0o00
	}

	millsStr := strconv.FormatFloat(mills, 'f', 0, strconv.IntSize)

	sec, ok := dateArray[0].(float64)
	if !ok {
		return time.Time{}
	}

	secStr := strconv.FormatFloat(sec, 'f', 0, strconv.IntSize)

	secs, err := strconv.ParseInt(secStr, 10, strconv.IntSize)
	if err != nil {
		return time.Time{}
	}

	nsecs, err := strconv.ParseInt(millsStr, 10, strconv.IntSize)
	if err != nil {
		return time.Unix(secs, 0)
	}

	return time.Unix(secs, nsecs)
}
