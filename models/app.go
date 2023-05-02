package models

import (
	"fmt"

	"github.com/bots-house/google-play-parser/internal/shared"
)

type App struct {
	AppID                    string
	URL                      string
	Title                    string
	Description              string
	DescriptionText          string
	Summary                  string
	Installs                 string
	MinInstalls              float64
	MaxInstalls              float64
	Currency                 string
	Price                    float64
	PriceText                string
	Free                     bool
	Score                    float64
	ScoreText                string
	Ratings                  float64
	Reviews                  float64
	Histogram                map[int]float64
	Available                bool
	OffersIAP                bool
	IAPRange                 string
	AndroidVersion           string
	AndroidVersionText       string
	Developer                string
	DeveloperID              string
	DeveloperEmail           string
	DeveloperWebsite         string
	DeveloperAddress         string
	PrivacyPolicy            string
	DeveloperInternalID      string
	Genre                    string
	GenreID                  string
	FamilyGenre              string
	FamilyGenreID            string
	Icon                     string
	HeaderImage              string
	Screenshots              []string
	Video                    string
	VideoImage               string
	PreviewVideo             string
	ContentRating            string
	ContentRatingDescription string
	AdSupported              bool
	Released                 string
	Updated                  float64
	Version                  string
	RecentChanges            string
	Comments                 []any
}

func (app App) Assign(rhs App) App {
	return shared.Assign(app, rhs)
}

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
