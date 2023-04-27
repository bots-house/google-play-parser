package googleplayscraper

import (
	"context"

	"github.com/bots-house/google-play-parser/models"
)

type Scrapper interface {
	Similar(ctx context.Context, opts models.ApplicationSpec) ([]models.App, error)
	App(ctx context.Context, opts models.ApplicationSpec) (models.App, error)
}
