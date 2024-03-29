package models

import (
	"fmt"

	"github.com/bots-house/google-play-parser/internal/shared"
)

type DeveloperSpec struct {
	DevID   string
	Lang    string
	Country string
	Count   int
	Full    bool
}

func (spec *DeveloperSpec) ensureNotNil() {
	*spec = shared.Assign(spec, &DeveloperSpec{Lang: "en", Country: "us", Count: 60})
}

func (spec *DeveloperSpec) Validate() error {
	spec.ensureNotNil()

	if spec.DevID == "" {
		return fmt.Errorf("devID required")
	}

	return nil
}
