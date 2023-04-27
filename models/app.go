package models

import "github.com/bots-house/google-play-parser/internal/shared"

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
	Comments                 []any // todo: no comments there
}

func (app App) Assign(rhs App) App {
	return shared.Assign(app, rhs)
}
