package gpp

import (
	"context"
	"net/http"

	"github.com/bots-house/google-play-parser/internal/scraper"
	"github.com/bots-house/google-play-parser/shared"
)

type collector struct {
	client shared.HTTPClient
}

var _ Scrapper = &collector{}

type CollectorOption func(*collector)

func New(opts ...CollectorOption) *collector {
	collector := &collector{
		client: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(collector)
	}

	return collector
}

func WithClient(client shared.HTTPClient) CollectorOption {
	return func(c *collector) {
		c.client = client
	}
}

func (collector collector) Similar(ctx context.Context, opts ApplicationSpec) ([]App, error) {
	apps, err := scraper.Similar(ctx, collector.client, opts.toInternal())
	if err != nil {
		return nil, err
	}

	return newApps(apps...), nil
}

func (collector collector) App(ctx context.Context, opts ApplicationSpec) (App, error) {
	app, err := scraper.App(ctx, collector.client, opts.toInternal())
	if err != nil {
		return App{}, err
	}

	return newFromInternal(&app), nil
}
