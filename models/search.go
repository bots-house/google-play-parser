package models

import (
	"fmt"

	"github.com/bots-house/google-play-parser/internal/shared"
)

type SearchSpec struct {
	Query   string
	Count   int
	Lang    string
	Country string
	Price   string
	Full    bool
}

func (spec *SearchSpec) ensureNotNil() {
	*spec = shared.Assign(spec, &SearchSpec{
		Lang:    "en",
		Country: "us",
		Count:   20,
		Price:   "all",
	})
}

func (spec *SearchSpec) Validate() error {
	spec.ensureNotNil()

	if spec.Query == "" {
		return fmt.Errorf("query required")
	}

	if spec.Count > 250 {
		return fmt.Errorf("apps count too large")
	}

	return nil
}
