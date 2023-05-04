package cli

import (
	"fmt"

	gpp "github.com/bots-house/google-play-parser"
	"github.com/urfave/cli/v2"
)

var developerCMD = &cli.Command{
	Name:    "developer",
	Aliases: []string{"d", "dev"},
	Flags: []cli.Flag{
		devIDFlag,
		countryFlag,
		langFlag,
		fullFlag,
	},
	Action: func(ctx *cli.Context) error {
		collector := gpp.New()

		if !ctx.IsSet(devIDFlag.Name) {
			return fmt.Errorf("dev id not set")
		}

		result, err := collector.Developer(ctx.Context, gpp.DeveloperSpec{
			DevID:   devIDFlag.Get(ctx),
			Lang:    langFlag.Get(ctx),
			Country: countryFlag.Get(ctx),
			Full:    fullFlag.Get(ctx),
		})
		if err != nil {
			return fmt.Errorf("developer method: %w", err)
		}

		return display(ctx, result)
	},
}
