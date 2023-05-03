package scraper

import (
	"strconv"
	"strings"

	"github.com/bots-house/google-play-parser/internal/ramda"
	"github.com/bots-house/google-play-parser/internal/shared"
)

const (
	baseURL        = "https://play.google.com"
	appsDetailsURL = "/store/apps/details"

	vary = "VARY"
)

type ParsedObject struct {
	Data        map[string][]any
	ServiceData map[string]Service
}

type Service struct {
	ID      string  `json:"id"`
	Ext     float64 `json:"ext,omitempty"`
	Request []any   `json:"request,omitempty"`
}

var clusterMapping = shared.ClusterMapping{
	Title: []any{21, 1, 0},
	URL:   []any{21, 1, 2, 4, 2},
	Apps:  []any{"ds:3", 0, 1, 0, 21, 0},
	Token: []any{"ds:3", 0, 1, 0, 21, 1, 3, 1},
}

var clusterSpec = shared.ParsedSpec{
	Clusters: shared.ParsedClustersSpec{
		Path:          []any{1, 1},
		UserServiceID: "ag2B9c",
	},
}

var appDetailsMapping = shared.Mapping{
	Title:     []any{"ds:5", 1, 2, 0, 0},
	Summary:   []any{"ds:5", 1, 2, 73, 0, 1},
	Score:     []any{"ds:5", 1, 2, 51, 0, 1},
	ScoreText: []any{"ds:5", 1, 2, 51, 0, 0},
	Price: shared.MappingWithFunc[float64, float64]{
		Path: []any{"ds:5", 1, 2, 57, 0, 0, 0, 0, 1, 0, 0},
		Fun: func(f float64) float64 {
			return f / 1000000
		},
	},
	PriceText: shared.MappingWithFunc[string, string]{
		Path: []any{"ds:5", 1, 2, 57, 0, 0, 0, 0, 1, 0, 2},
		Fun: func(f string) string {
			if f == "" {
				return "Free"
			}

			return f
		},
	},
	Free: shared.MappingWithFunc[float64, bool]{
		Path: []any{"ds:5", 1, 2, 57, 0, 0, 0, 0, 1, 0, 0},
		Fun: func(f float64) bool {
			return f == 0
		},
	},
	Currency: []any{"ds:5", 1, 2, 57, 0, 0, 0, 0, 1, 0, 1},
	Icon:     []any{"ds:5", 1, 2, 95, 0, 3, 2},
	Description: shared.MappingWithFunc[[]any, string]{
		Path: []any{"ds:5", 1, 2},
		Fun: func(a []any) string {
			description := parseDescription(a)

			return strings.ReplaceAll(description, "<br>", "\r\n")
		},
	},
	DescriptionText: shared.MappingWithFunc[[]any, string]{
		Path: []any{"ds:5", 1, 2},
		Fun:  parseDescription,
	},
	Installs:    []any{"ds:5", 1, 2, 13, 0},
	MinInstalls: []any{"ds:5", 1, 2, 13, 1},
	MaxInstalls: []any{"ds:5", 1, 2, 13, 2},
	Ratings:     []any{"ds:5", 1, 2, 51, 2, 1},
	Reviews:     []any{"ds:5", 1, 2, 51, 3, 1},
	Histogram: shared.MappingWithFunc[[]any, map[int]float64]{
		Path: []any{"ds:5", 1, 2, 51, 1},
		Fun: func(a []any) map[int]float64 {
			histogram := map[int]float64{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}
			if len(a) == 0 {
				return histogram
			}

			for idx, path := range [][]any{
				{1, 1},
				{2, 1},
				{3, 1},
				{4, 1},
				{5, 1},
			} {
				value, ok := ramda.Path(path, a).(float64)
				if !ok {
					continue
				}

				histogram[idx+1] = value
			}

			return histogram
		},
	},
	Available: shared.MappingWithFunc[float64, bool]{
		Path: []any{"ds:5", 1, 2, 18, 0},
		Fun:  func(f float64) bool { return f > 0 },
	},
	OffersIAP: shared.MappingWithFunc[string, bool]{
		Path: []any{"ds:5", 1, 2, 19, 0},
		Fun:  func(s string) bool { return len(s) > 0 },
	},
	IAPRange: []any{"ds:5", 1, 2, 19, 0},
	AndroidVersion: shared.MappingWithFunc[string, string]{
		Path: []any{"ds:5", 1, 2, 140, 1, 1, 0, 0, 1},
		Fun: func(s string) string {
			if s == "" {
				return vary
			}

			number := strings.Split(s, " ")[0]

			_, err := strconv.ParseFloat(number, 64)
			if err != nil {
				return vary
			}

			return number
		},
	},
	AndroidVersionText: shared.MappingWithFunc[string, string]{
		Path: []any{"ds:5", 1, 2, 140, 1, 1, 0, 0, 1},
		Fun: func(s string) string {
			if s == "" {
				return "Varies with device"
			}

			return s
		},
	},
	Developer: []any{"ds:5", 1, 2, 68, 1, 4, 2},
	DeveloperID: shared.MappingWithFunc[string, string]{
		Path: []any{"ds:5", 1, 2, 68, 1, 4, 2},
		Fun: func(s string) string {
			return strings.Split(s, "id=")[1]
		},
	},
	DeveloperEmail:   []any{"ds:5", 1, 2, 69, 1, 0},
	DeveloperWebsite: []any{"ds:5", 1, 2, 69, 0, 5, 2},
	DeveloperAddress: []any{"ds:5", 1, 2, 69, 2, 0},
	PrivacyPolicy:    []any{"ds:5", 1, 2, 99, 0, 5, 2},
	DeveloperInternalID: shared.MappingWithFunc[string, string]{
		Path: []any{"ds:5", 1, 2, 68, 1, 4, 2},
		Fun: func(s string) string {
			return strings.Split(s, "id=")[1]
		},
	},
	Genre:         []any{"ds:5", 1, 2, 79, 0, 0, 0},
	GenreID:       []any{"ds:5", 1, 2, 79, 0, 0, 2},
	FamilyGenre:   []any{"ds:5", 0, 12, 13, 1, 0},
	FamilyGenreID: []any{"ds:5", 0, 12, 13, 1, 2},
	HeaderImage:   []any{"ds:5", 1, 2, 95, 0, 3, 2},
	Screenshots: shared.MappingWithFunc[[]any, []string]{
		Path: []any{"ds:5", 1, 2, 78, 0},
		Fun: func(a []any) []string {
			result := make([]string, 0, len(a))

			for _, screenshot := range a {
				raw, ok := ramda.Path([]any{3, 2}, screenshot).(string)
				if !ok {
					continue
				}

				result = append(result, raw)
			}

			return result
		},
	},
	Video:                    []any{"ds:5", 1, 2, 100, 0, 0, 3, 2},
	VideoImage:               []any{"ds:5", 1, 2, 100, 1, 0, 3, 2},
	PreviewVideo:             []any{"ds:5", 1, 2, 100, 1, 2, 0, 2},
	ContentRating:            []any{"ds:5", 1, 2, 9, 0},
	ContentRatingDescription: []any{"ds:5", 1, 2, 9, 2, 1},
	AdSupported: shared.MappingWithFunc[[]any, bool]{
		Path: []any{"ds:5", 1, 2, 48},
		Fun:  func(a []any) bool { return len(a) > 0 },
	},
	Released: []any{"ds:5", 1, 2, 10, 0},
	Updated: shared.MappingWithFunc[float64, float64]{
		Path: []any{"ds:5", 1, 2, 145, 0, 1, 0},
		Fun: func(i float64) float64 {
			return i * 1000
		},
	},
	Version: shared.MappingWithFunc[string, string]{
		Path: []any{"ds:5", 1, 2, 140, 0, 0, 0},
		Fun: func(s string) string {
			if s == "" {
				return vary
			}

			return s
		},
	},
	RecentChanges: []any{"ds:5", 1, 2, 144, 1, 1},
	Comments:      []any{"ds:9", 0},
}
