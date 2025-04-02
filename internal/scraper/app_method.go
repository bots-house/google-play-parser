package scraper

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/bots-house/google-play-parser/internal/parser"
	"github.com/bots-house/google-play-parser/models"
	sh "github.com/bots-house/google-play-parser/shared"
)

func App(ctx context.Context, client sh.HTTPClient, spec models.ApplicationSpec) (models.App, error) {
	if err := spec.Validate(); err != nil {
		return models.App{}, fmt.Errorf("validation: %w", err)
	}

	appURL := getURL(appsDetailsURL)

	body, requestURL, err := request(ctx, client, requestSpec{
		url: appURL,
		params: url.Values{
			"id": []string{spec.AppID},
			"gl": []string{spec.Country},
			"hl": []string{spec.Lang},
		},
	})
	if err != nil {
		return models.App{}, err
	}

	parsed, err := parser.Parse(body)
	if err != nil {
		return models.App{}, fmt.Errorf("parse: %w", err)
	}

	app, ok := parser.Extract[models.App](parsed.Data, &appDetailsMapping)
	if !ok {
		return models.App{}, fmt.Errorf("no app details found")
	}

	app.Developer = strings.Split(app.Developer, "id=")[1]
	app = checkDeveloperName(ctx, client, app)
	app.Unquote()

	return app.Assign(&models.App{AppID: spec.AppID, URL: requestURL}), nil
}

func checkDeveloperName(ctx context.Context, client sh.HTTPClient, app models.App) models.App {
	name := app.Developer
	if _, err := strconv.ParseInt(name, 10, strconv.IntSize); err != nil {
		return app
	}

	devApps, err := Developer(ctx, client, models.DeveloperSpec{DevID: app.DeveloperID})
	if err != nil {
		return app
	}

	if len(devApps) == 0 {
		return app
	}

	app.Developer = devApps[0].Developer

	return app
}
