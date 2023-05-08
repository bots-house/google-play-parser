package gpp

import (
	"context"
)

type Scrapper interface {
	Similar(context.Context, ApplicationSpec) ([]App, error)
	App(context.Context, ApplicationSpec) (App, error)
	List(context.Context, ListSpec) ([]App, error)
	Developer(context.Context, DeveloperSpec) ([]App, error)
	Search(context.Context, SearchSpec) ([]App, error)
	DataSafety(context.Context, ApplicationSpec) (DataSafety, error)
}

// TODO: add proces over pages for methods
