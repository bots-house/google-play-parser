package cli

import (
	"fmt"

	gpp "github.com/bots-house/google-play-parser"
	"github.com/urfave/cli/v2"
)

var suggestCMD = &cli.Command{
	Name: "suggest",
	Flags: []cli.Flag{
		queryFlag,
		countryFlag,
		langFlag,
	},
	Action: func(ctx *cli.Context) error {
		collector := gpp.New()

		if !ctx.IsSet(queryFlag.Name) {
			return fmt.Errorf("query not set")
		}

		result, err := collector.Suggest(ctx.Context, gpp.SearchSpec{
			Query:   queryFlag.Get(ctx),
			Lang:    langFlag.Get(ctx),
			Country: countryFlag.Get(ctx),
		})
		if err != nil {
			return fmt.Errorf("suggest method: %w", err)
		}

		return display(ctx, result)
	},
}
