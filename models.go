package gpp

import (
	"github.com/bots-house/google-play-parser/internal/shared"
	"github.com/bots-house/google-play-parser/models"
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
	Comments                 []any           `json:"comments,omitempty"`
}

func newFromInternal(app *models.App) App {
	return App{
		AppID:                    app.AppID,
		URL:                      app.URL,
		Title:                    app.Title,
		Description:              app.Description,
		DescriptionText:          app.DescriptionText,
		Summary:                  app.Summary,
		Installs:                 app.Installs,
		MinInstalls:              app.MinInstalls,
		MaxInstalls:              app.MaxInstalls,
		Currency:                 app.Currency,
		Price:                    app.Price,
		PriceText:                app.PriceText,
		Free:                     app.Free,
		Score:                    app.Score,
		ScoreText:                app.ScoreText,
		Ratings:                  app.Ratings,
		Reviews:                  app.Reviews,
		Histogram:                app.Histogram,
		Available:                app.Available,
		OffersIAP:                app.OffersIAP,
		IAPRange:                 app.IAPRange,
		AndroidVersion:           app.AndroidVersion,
		AndroidVersionText:       app.AndroidVersionText,
		Developer:                app.Developer,
		DeveloperID:              app.DeveloperID,
		DeveloperEmail:           app.DeveloperEmail,
		DeveloperWebsite:         app.DeveloperWebsite,
		DeveloperAddress:         app.DeveloperAddress,
		Genre:                    app.Genre,
		GenreID:                  app.GenreID,
		FamilyGenre:              app.FamilyGenre,
		FamilyGenreID:            app.FamilyGenreID,
		Icon:                     app.Icon,
		HeaderImage:              app.HeaderImage,
		Screenshots:              app.Screenshots,
		Video:                    app.Video,
		VideoImage:               app.VideoImage,
		PreviewVideo:             app.PreviewVideo,
		ContentRating:            app.ContentRating,
		PrivacyPolicy:            app.PrivacyPolicy,
		ContentRatingDescription: app.ContentRatingDescription,
		AdSupported:              app.AdSupported,
		Released:                 app.Released,
		Updated:                  app.Updated,
		Version:                  app.Version,
		RecentChanges:            app.RecentChanges,
		Comments:                 app.Comments,
	}
}

func newApps(apps ...models.App) []App {
	result := make([]App, 0, len(apps))

	for _, app := range apps {
		result = append(result, newFromInternal(&app))
	}

	return result
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

func (spec ApplicationSpec) toInternal() models.ApplicationSpec {
	return models.ApplicationSpec{
		AppID:   spec.AppID,
		Country: spec.Country,
		Lang:    spec.Lang,
		Full:    spec.Full,
	}
}
