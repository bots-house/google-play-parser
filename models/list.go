package models

import (
	"fmt"

	"github.com/bots-house/google-play-parser/internal/shared"
)

var (
	listCollections = []string{"TOP_FREE", "TOP_PAID", "GROSSING"}
	listAges        = []string{"AGE_RANGE1", "AGE_RANGE2", "AGE_RANGE3"}
	listCategories  = []string{
		"APPLICATION",
		"ANDROID_WEAR",
		"ART_AND_DESIGN",
		"AUTO_AND_VEHICLES",
		"BEAUTY",
		"BOOKS_AND_REFERENCE",
		"BUSINESS",
		"COMICS",
		"COMMUNICATION",
		"DATING",
		"EDUCATION",
		"ENTERTAINMENT",
		"EVENTS",
		"FINANCE",
		"FOOD_AND_DRINK",
		"HEALTH_AND_FITNESS",
		"HOUSE_AND_HOME",
		"LIBRARIES_AND_DEMO",
		"LIFESTYLE",
		"MAPS_AND_NAVIGATION",
		"MEDICAL",
		"MUSIC_AND_AUDIO",
		"NEWS_AND_MAGAZINES",
		"PARENTING",
		"PERSONALIZATION",
		"PHOTOGRAPHY",
		"PRODUCTIVITY",
		"SHOPPING",
		"SOCIAL",
		"SPORTS",
		"TOOLS",
		"TRAVEL_AND_LOCAL",
		"VIDEO_PLAYERS",
		"WATCH_FACE",
		"WEATHER",
		"GAME",
		"GAME_ACTION",
		"GAME_ADVENTURE",
		"GAME_ARCADE",
		"GAME_BOARD",
		"GAME_CARD",
		"GAME_CASINO",
		"GAME_CASUAL",
		"GAME_EDUCATIONAL",
		"GAME_MUSIC",
		"GAME_PUZZLE",
		"GAME_RACING",
		"GAME_ROLE_PLAYING",
		"GAME_SIMULATION",
		"GAME_SPORTS",
		"GAME_STRATEGY",
		"GAME_TRIVIA",
		"GAME_WORD",
		"FAMILY",
	}
)

const (
	defaultListCount = 100
)

func GetDefaultListCount() int {
	return defaultListCount
}

type ListSpec struct {
	Count      int
	Age        string
	Lang       string
	Country    string
	Category   string
	Collection string
	Full       bool
}

func (spec *ListSpec) ensureNotNil() {
	*spec = shared.Assign(
		spec,
		&ListSpec{Lang: "en", Country: "us", Category: listCategories[0], Collection: listCollections[0], Count: defaultListCount},
	)
}

func (spec *ListSpec) Validate() error {
	spec.ensureNotNil()

	if spec.Count < 1 {
		return fmt.Errorf("count must be greater then zero")
	}

	if !shared.In(spec.Category, listCategories...) {
		return fmt.Errorf("invalid category")
	}

	if !shared.In(spec.Collection, listCollections...) {
		return fmt.Errorf("invalid collection")
	}

	if spec.Age != "" && !shared.In(spec.Age, listAges...) {
		return fmt.Errorf("invalid age")
	}

	return nil
}
