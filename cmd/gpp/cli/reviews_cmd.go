package cli

import (
	"fmt"
	"strings"

	gpp "github.com/bots-house/google-play-parser"
	"github.com/urfave/cli/v2"
)

var reviewsCMD = &cli.Command{
	Name: "reviews",
	Flags: []cli.Flag{
		appIDFlag,
		countryFlag,
		langFlag,
		countFlag,
		sortFlag,
	},
	Action: func(ctx *cli.Context) error {
		collector := gpp.New()

		if !ctx.IsSet(appIDFlag.Name) {
			return fmt.Errorf("app id not set")
		}

		result, err := collector.Reviews(ctx.Context, gpp.ReviewsSpec{
			AppID:   appIDFlag.Get(ctx),
			Lang:    langFlag.Get(ctx),
			Country: countryFlag.Get(ctx),
			Count:   countFlag.Get(ctx),
			Sort:    parseSortFlag(sortFlag.Get(ctx)),
		})
		if err != nil {
			return fmt.Errorf("app method: %w", err)
		}

		return display(ctx, result)
	},
}

func parseSortFlag(sort string) gpp.ReviewsSort {
	switch strings.ToLower(sort) {
	case "helpfulness":
		return 1
	case "newest", "":
		return 2
	case "rating":
		return 3
	default:
		return -1
	}
}
