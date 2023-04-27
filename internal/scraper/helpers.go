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

type requestSpec struct {
	method string
	url    string
	params *url.Values
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

func request(ctx context.Context, client shared.HTTPClient, spec requestSpec) ([]byte, error) {
	spec.ensureNotNil()

	if err := spec.validate(); err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, spec.method, spec.url, nil)
	if err != nil {
		return nil, fmt.Errorf("prepare request: %w", err)
	}

	if spec.params != nil {
		request.URL.RawQuery = spec.params.Encode()
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response: status: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return body, nil
}

func parseDescription(description []any) string {
	translation, _ := ramda.Path([]any{12, 0, 0, 1}, description).(string)
	original, _ := ramda.Path([]any{72, 0, 1}, description).(string)

	if translation == "" {
		return original
	}

	return translation
}
