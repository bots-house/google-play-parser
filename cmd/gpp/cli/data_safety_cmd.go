package cli

import (
	"fmt"

	gpp "github.com/bots-house/google-play-parser"
	"github.com/urfave/cli/v2"
)

var dataSafetyCMD = &cli.Command{
	Name: "data-safety",
	Flags: []cli.Flag{
		appIDFlag,
		countryFlag,
		langFlag,
		fullFlag,
	},
	Action: func(ctx *cli.Context) error {
		collector := gpp.New()

		if !ctx.IsSet(appIDFlag.Name) {
			return fmt.Errorf("app id not set")
		}

		result, err := collector.DataSafety(ctx.Context, gpp.ApplicationSpec{
			AppID: appIDFlag.Get(ctx),
			Lang:  langFlag.Get(ctx),
		})
		if err != nil {
			return fmt.Errorf("data_safety method: %w", err)
		}

		return display(ctx, result)
	},
}
