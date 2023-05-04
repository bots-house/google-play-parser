package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/bots-house/google-play-parser/internal/ramda"
	"github.com/bots-house/google-play-parser/shared"
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

func request(ctx context.Context, client shared.HTTPClient, spec requestSpec) (body []byte, rawURL string, err error) {
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
