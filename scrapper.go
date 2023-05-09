package gpp

import (
	"context"
	"net/http"

	"github.com/bots-house/google-play-parser/internal/scraper"
	"github.com/bots-house/google-play-parser/internal/shared"
	"github.com/bots-house/google-play-parser/models"
	sh "github.com/bots-house/google-play-parser/shared"
)

type collector struct {
	client sh.HTTPClient
}

var _ Scrapper = &collector{}

type CollectorOption func(*collector)

func New(opts ...CollectorOption) Scrapper {
	collector := &collector{
		client: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(collector)
	}

	return collector
}

func WithClient(client sh.HTTPClient) CollectorOption {
	return func(c *collector) {
		c.client = client
	}
}

func (collector collector) Similar(ctx context.Context, spec ApplicationSpec) ([]App, error) {
	apps, err := scraper.Similar(ctx, collector.client, spec.toInternal())
	if err != nil {
		return nil, err
	}

	return newApps(apps...), nil
}

func (collector collector) App(ctx context.Context, spec ApplicationSpec) (App, error) {
	app, err := scraper.App(ctx, collector.client, spec.toInternal())
	if err != nil {
		return App{}, err
	}

	return newFromInternal(&app), nil
}

func (collector collector) List(ctx context.Context, spec ListSpec) ([]App, error) {
	apps, err := scraper.List(ctx, collector.client, spec.toInternal())
	if err != nil {
		return nil, err
	}

	return newApps(apps...), nil
}

func (collector collector) Developer(ctx context.Context, spec DeveloperSpec) ([]App, error) {
	apps, err := scraper.Developer(ctx, collector.client, spec.toInternal())
	if err != nil {
		return nil, err
	}

	return newApps(apps...), nil
}

func (collector collector) Search(ctx context.Context, spec SearchSpec) ([]App, error) {
	apps, err := scraper.Search(ctx, collector.client, models.SearchSpec(spec))
	if err != nil {
		return nil, err
	}

	return newApps(apps...), nil
}

func (collector collector) DataSafety(ctx context.Context, spec ApplicationSpec) (DataSafety, error) {
	dataSafety, err := scraper.DataSafety(ctx, collector.client, spec.toInternal())
	if err != nil {
		return DataSafety{}, err
	}

	return DataSafety(dataSafety), nil
}

func (collector collector) Permissions(ctx context.Context, spec ApplicationSpec) ([]Permission, error) {
	perms, err := scraper.Permissions(ctx, collector.client, spec.toInternal())
	if err != nil {
		return nil, err
	}

	return shared.Map(perms, func(perm models.Permission) Permission {
		return Permission(perm)
	}), nil
}

func (collector collector) Suggest(ctx context.Context, spec SearchSpec) ([]string, error) {
	return scraper.Suggest(ctx, collector.client, models.SearchSpec(spec))
}
