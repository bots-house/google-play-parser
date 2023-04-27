package models

import "fmt"

type ApplicationSpec struct {
	AppID   string
	Lang    string
	Country string
	Full    bool
}

var defaultSimilarOpts = ApplicationSpec{
	Lang:    "en",
	Country: "us",
}

func (opts *ApplicationSpec) EnsureNotNil() {
	if opts.Lang == "" {
		opts.Lang = defaultSimilarOpts.Lang
	}

	if opts.Country == "" {
		opts.Country = defaultSimilarOpts.Country
	}
}

func (opts ApplicationSpec) Validate() error {
	if opts.AppID == "" {
		return fmt.Errorf("appID required")
	}

	return nil
}
