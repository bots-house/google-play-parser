package models

import (
	"fmt"
	"net/url"
	"strings"

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
	InAppPurchase            bool
}

func (app *App) Assign(rhs *App) App {
	return shared.Assign(app, rhs)
}

func (app *App) Unquote() {
	developer, err := url.QueryUnescape(app.Developer)
	if err != nil {
		developer = app.Developer
	}

	developerID, err := url.QueryUnescape(app.DeveloperID)
	if err != nil {
		developerID = app.DeveloperID
	}

	developer = strings.ReplaceAll(developer, "+", " ")

	app.Developer = developer
	app.DeveloperID = developerID
}

type ApplicationSpec struct {
	AppID   string
	Lang    string
	Country string
	Count   int
	Full    bool
}

var defaultSimilarSpec = ApplicationSpec{
	Lang:    "en",
	Country: "us",
}

func (spec *ApplicationSpec) ensureNotNil() {
	if spec.Lang == "" {
		spec.Lang = defaultSimilarSpec.Lang
	}

	if spec.Country == "" {
		spec.Country = defaultSimilarSpec.Country
	}
}

func (spec *ApplicationSpec) Validate() error {
	spec.ensureNotNil()

	if spec.AppID == "" {
		return fmt.Errorf("appID required")
	}

	return nil
}
