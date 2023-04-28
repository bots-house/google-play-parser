package models

import (
	"fmt"

	"github.com/bots-house/google-play-parser/internal/shared"
)

type App struct {
	AppID                    string          `json:"app_id"`
	URL                      string          `json:"url"`
	Title                    string          `json:"title"`
	Description              string          `json:"description,omitempty"`
	DescriptionText          string          `json:"description_text,omitempty"`
	Summary                  string          `json:"summary"`
	Installs                 string          `json:"installs,omitempty"`
	MinInstalls              float64         `json:"min_installs,omitempty"`
	MaxInstalls              float64         `json:"max_installs,omitempty"`
	Currency                 string          `json:"currency"`
	Price                    float64         `json:"price,omitempty"`
	PriceText                string          `json:"price_text,omitempty"`
	Free                     bool            `json:"free"`
	Score                    float64         `json:"score"`
	ScoreText                string          `json:"score_text"`
	Ratings                  float64         `json:"ratings,omitempty"`
	Reviews                  float64         `json:"reviews,omitempty"`
	Histogram                map[int]float64 `json:"histogram,omitempty"`
	Available                bool            `json:"available,omitempty"`
	OffersIAP                bool            `json:"offers_iap,omitempty"`
	IAPRange                 string          `json:"iap_range,omitempty"`
	AndroidVersion           string          `json:"android_version,omitempty"`
	AndroidVersionText       string          `json:"android_version_text,omitempty"`
	Developer                string          `json:"developer"`
	DeveloperID              string          `json:"developer_id,omitempty"`
	DeveloperEmail           string          `json:"developer_email,omitempty"`
	DeveloperWebsite         string          `json:"developer_website,omitempty"`
	DeveloperAddress         string          `json:"developer_address,omitempty"`
	PrivacyPolicy            string          `json:"privacy_policy,omitempty"`
	DeveloperInternalID      string          `json:"developer_internal_id,omitempty"`
	Genre                    string          `json:"genre,omitempty"`
	GenreID                  string          `json:"genre_id,omitempty"`
	FamilyGenre              string          `json:"family_genre,omitempty"`
	FamilyGenreID            string          `json:"family_genre_id,omitempty"`
	Icon                     string          `json:"icon"`
	HeaderImage              string          `json:"header_image,omitempty"`
	Screenshots              []string        `json:"screenshots,omitempty"`
	Video                    string          `json:"video,omitempty"`
	VideoImage               string          `json:"video_image,omitempty"`
	PreviewVideo             string          `json:"preview_video,omitempty"`
	ContentRating            string          `json:"content_rating,omitempty"`
	ContentRatingDescription string          `json:"content_rating_description,omitempty"`
	AdSupported              bool            `json:"ad_supported,omitempty"`
	Released                 string          `json:"released,omitempty"`
	Updated                  float64         `json:"updated,omitempty"`
	Version                  string          `json:"version,omitempty"`
	RecentChanges            string          `json:"recent_changes,omitempty"`
	Comments                 []any           `json:"comments,omitempty"` // todo: no comments there
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
