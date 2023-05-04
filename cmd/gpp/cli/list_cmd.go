package cli

import (
	"fmt"

	gpp "github.com/bots-house/google-play-parser"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var listCMD = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Flags: []cli.Flag{
		appIDFlag,
		countryFlag,
		langFlag,
		countFlag,
		categoryFlag,
		collectionFlag,
		ageFlag,
		fullFlag,
	},
	Action: func(ctx *cli.Context) error {
		collector := gpp.New()

		result, err := collector.List(ctx.Context, gpp.ListSpec{
			Count:      countFlag.Get(ctx),
			Lang:       langFlag.Get(ctx),
			Country:    countryFlag.Get(ctx),
			Collection: collectionFlag.Get(ctx),
			Category:   categoryFlag.Get(ctx),
			Age:        ageFlag.Get(ctx),
			Full:       fullFlag.Get(ctx),
		})
		if err != nil {
			return fmt.Errorf("list method: %w", err)
		}

		log.Ctx(ctx.Context).Debug().Msgf("similar apps founded: %d", len(result))

		return display(ctx, result)
	},
}
