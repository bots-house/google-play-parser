package gpp

import (
	"context"
)

type Scrapper interface {
	Similar(context.Context, ApplicationSpec) ([]App, error)
	App(context.Context, ApplicationSpec) (App, error)
	List(context.Context, ListSpec) ([]App, error)
	Developer(context.Context, DeveloperSpec) ([]App, error)
}
