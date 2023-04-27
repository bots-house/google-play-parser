package scraper

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/bots-house/google-play-parser/internal/parser"
	"github.com/bots-house/google-play-parser/models"
	sh "github.com/bots-house/google-play-parser/shared"
)

func App(ctx context.Context, client sh.HTTPClient, opts models.ApplicationSpec) (models.App, error) {
	opts.EnsureNotNil()

	if err := opts.Validate(); err != nil {
		return models.App{}, err
	}

	appURL := getURL(appsDetailsURL)

	body, requestURL, err := request(ctx, client, requestSpec{
		url: appURL,
		params: &url.Values{
			"id": []string{opts.AppID},
			"gl": []string{opts.Country},
			"hl": []string{opts.Lang},
		},
	})
	if err != nil {
		return models.App{}, err
	}

	parsed, err := parser.Parse(body)
	if err != nil {
		return models.App{}, fmt.Errorf("parse: %w", err)
	}

	app, ok := parser.Extract(parsed.Data, appDetailsMapping)
	if !ok {
		return models.App{}, fmt.Errorf("no app details found")
	}

	app.Developer = strings.Split(app.Developer, "id=")[1]

	return app.Assign(models.App{AppID: opts.AppID, URL: requestURL}), nil
}
