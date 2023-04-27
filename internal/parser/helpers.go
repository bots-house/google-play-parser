package parser

import (
	"fmt"
	"regexp"

	"go.uber.org/multierr"
)

func matches(data []byte, patterns ...*regexp.Regexp) (matches error) {
	errs := make([]error, len(patterns))

	for _, pattern := range patterns {
		if !pattern.Match(data) {
			errs = append(errs, fmt.Errorf("did`nt match pattern: [%s]", pattern))
		}
	}

	if len(errs) > 0 {
		matches = multierr.Combine(errs...)
	}

	return
}
