package gpp

import (
	"context"
	"net/http"

	"github.com/bots-house/google-play-parser/internal/scraper"
	"github.com/bots-house/google-play-parser/models"
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

func (collector collector) Similar(ctx context.Context, opts models.ApplicationSpec) ([]models.App, error) {
	return scraper.Similar(ctx, collector.client, opts)
}

func (collector collector) App(ctx context.Context, opts models.ApplicationSpec) (models.App, error) {
	return scraper.App(ctx, collector.client, opts)
}
