package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/bots-house/google-play-parser/internal/shared"
)

type ReviewsSpec struct {
	AppID       string
	Lang        string
	Country     string
	Count       int
	Sort        int
	Paginated   bool
	NextToken   string
	RequestType string
}

func (spec *ReviewsSpec) ensureNotNil() {
	*spec = shared.Assign(
		spec,
		&ReviewsSpec{Sort: 2, Lang: "en", Country: "us", Count: 150},
	)

	if spec.NextToken == "" {
		spec.RequestType = "initial"
		return
	}

	spec.RequestType = "paginated"
}

func (spec *ReviewsSpec) Validate() error {
	spec.ensureNotNil()

	if spec.AppID == "" {
		return fmt.Errorf("app id required")
	}

	if !shared.In(spec.Sort, 1, 2, 3) {
		return fmt.Errorf("sort must be in [1,2,3]")
	}

	return nil
}

type Review struct {
	ID        uuid.UUID        `json:"id"`
	URL       string           `json:"url"`
	Title     string           `json:"title,omitempty"`
	Summary   string           `json:"summary"`
	Score     float64          `json:"score"`
	ScoreText string           `json:"score_text"`
	UserImage string           `json:"user_image"`
	UserName  string           `json:"user_name"`
	Version   string           `json:"version"`
	Date      time.Time        `json:"date"`
	ReplyText string           `json:"reply_text,omitempty"`
	ReplyDate time.Time        `json:"reply_date,omitempty"`
	Criteria  []ReviewCriteria `json:"criteria"`
	TumbsUp   float64          `json:"tumbs_up"`
}

type ReviewCriteria struct {
	Criteria string  `json:"criteria"`
	Rating   float64 `json:"rating,omitempty"`
}
