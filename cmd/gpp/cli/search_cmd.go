package cli

import (
	"fmt"

	gpp "github.com/bots-house/google-play-parser"
	"github.com/urfave/cli/v2"
)

var searchCMD = &cli.Command{
	Name: "search",
	Flags: []cli.Flag{
		queryFlag,
		countryFlag,
		countFlag,
		langFlag,
		priceFlag,
		fullFlag,
	},
	Action: func(ctx *cli.Context) error {
		collector := gpp.New()

		if !ctx.IsSet(queryFlag.Name) {
			return fmt.Errorf("query not set")
		}

		result, err := collector.Search(ctx.Context, gpp.SearchSpec{
			Query:   queryFlag.Get(ctx),
			Lang:    langFlag.Get(ctx),
			Country: countryFlag.Get(ctx),
			Count:   countFlag.Get(ctx),
			Price:   priceFlag.Get(ctx),
			Full:    fullFlag.Get(ctx),
		})
		if err != nil {
			return fmt.Errorf("search method: %w", err)
		}

		return display(ctx, result)
	},
}
