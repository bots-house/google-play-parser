package cli

import (
	"fmt"

	gpp "github.com/bots-house/google-play-parser"
	"github.com/urfave/cli/v2"
)

var appCMD = &cli.Command{
	Name: "app",
	Flags: []cli.Flag{
		appIDFlag,
		countryFlag,
		langFlag,
	},
	Action: func(ctx *cli.Context) error {
		collector := gpp.New()

		if !ctx.IsSet(appIDFlag.Name) {
			return fmt.Errorf("app id not set")
		}

		result, err := collector.App(ctx.Context, gpp.ApplicationSpec{
			AppID:   appIDFlag.Get(ctx),
			Lang:    langFlag.Get(ctx),
			Country: countryFlag.Get(ctx),
		})
		if err != nil {
			return fmt.Errorf("similar method: %w", err)
		}

		return display(ctx, result)
	},
}
