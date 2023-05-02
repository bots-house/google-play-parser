package gpp

import (
	"context"
)

type Scrapper interface {
	Similar(ctx context.Context, opts ApplicationSpec) ([]App, error)
	App(ctx context.Context, opts ApplicationSpec) (App, error)
}
