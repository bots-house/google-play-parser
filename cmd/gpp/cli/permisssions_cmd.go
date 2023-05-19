package cli

import (
	"fmt"

	gpp "github.com/bots-house/google-play-parser"
	"github.com/urfave/cli/v2"
)

var permissionsCMD = &cli.Command{
	Name:    "permissions",
	Aliases: []string{"perm"},
	Flags: []cli.Flag{
		appIDFlag,
		langFlag,
		fullFlag,
	},
	Action: func(ctx *cli.Context) error {
		collector := gpp.New()

		if !ctx.IsSet(appIDFlag.Name) {
			return fmt.Errorf("app id not set")
		}

		result, err := collector.Permissions(ctx.Context, gpp.ApplicationSpec{
			AppID: appIDFlag.Get(ctx),
			Lang:  langFlag.Get(ctx),
			Full:  fullFlag.Get(ctx),
		})
		if err != nil {
			return fmt.Errorf("permissions method: %w", err)
		}

		return display(ctx, result)
	},
}
