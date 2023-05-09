package shared

import (
	"time"
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

type ParsedSpec struct {
	Clusters ParsedClustersSpec
}

type ClusterMapping struct {
	Title    []any
	URL      []any
	Apps     []any
	Token    []any
	Sections []any
}

type ParsedClustersSpec struct {
	Path          []any
	UserServiceID string
}

type AppMapping struct {
	AppID                    []any
	URL                      MappingWithFunc[string, string]
	Title                    []any
	Description              MappingWithFunc[[]any, string]
	DescriptionText          MappingWithFunc[[]any, string]
	Summary                  []any
	Installs                 []any
	MinInstalls              []any
	MaxInstalls              []any
	Currency                 []any
	Price                    MappingWithFunc[float64, float64]
	PriceText                MappingWithFunc[string, string]
	Free                     MappingWithFunc[float64, bool]
	ScoreText                []any
	Score                    []any
	Ratings                  []any
	Reviews                  []any
	Histogram                MappingWithFunc[[]any, map[int]float64]
	Available                MappingWithFunc[float64, bool]
	OffersIAP                MappingWithFunc[string, bool]
	IAPRange                 []any
	AndroidVersion           MappingWithFunc[string, string]
	AndroidVersionText       MappingWithFunc[string, string]
	Developer                []any
	DeveloperID              MappingWithFunc[string, string]
	DeveloperEmail           []any
	DeveloperWebsite         []any
	DeveloperAddress         []any
	PrivacyPolicy            []any
	DeveloperInternalID      MappingWithFunc[string, string]
	Genre                    []any
	GenreID                  []any
	FamilyGenre              []any
	FamilyGenreID            []any
	Icon                     []any
	HeaderImage              []any
	Screenshots              MappingWithFunc[[]any, []string]
	Video                    []any
	VideoImage               []any
	PreviewVideo             []any
	ContentRating            []any
	ContentRatingDescription []any
	AdSupported              MappingWithFunc[[]any, bool]
	Released                 []any
	Updated                  MappingWithFunc[float64, float64]
	Version                  MappingWithFunc[string, string]
	RecentChanges            []any
	Comments                 []any // TODO: no comments by this path
}

type MappingWithFunc[I, O any] struct {
	Path []any
	Fun  func(I) O
}

type DataSafetyMapping struct {
	SharedData        MappingWithFunc[any, []map[string]any]
	CollectedData     MappingWithFunc[any, []map[string]any]
	SecurityPractices MappingWithFunc[any, []map[string]any]
	PrivacyPolicyURL  []any
}

type ReviewsMapping struct {
	ID        []any
	UserName  []any
	UserImage []any
	Date      MappingWithFunc[[]any, time.Time]
	Score     []any
	ScoreText MappingWithFunc[float64, string]
	URL       MappingWithFunc[string, string]
	ReplyDate MappingWithFunc[[]any, time.Time]
	Summary   []any
	ReplyText []any
	Version   []any
	TumbsUp   []any
	Criteria  MappingWithFunc[[]any, map[string]float64]
}

type Mapping interface {
	*AppMapping | DataSafetyMapping | ReviewsMapping
}
