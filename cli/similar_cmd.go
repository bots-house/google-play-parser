package cli

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	googleplayscraper "github.com/bots-house/google-play-parser"
	"github.com/bots-house/google-play-parser/models"
)

var similarCMD = &cli.Command{
	Name:    "similar",
	Aliases: []string{"s"},
	Flags: []cli.Flag{
		appIDFlag,
		countryFlag,
		langFlag,
	},
	Action: func(ctx *cli.Context) error {
		collector := googleplayscraper.New()

		if !ctx.IsSet(appIDFlag.Name) {
			return fmt.Errorf("app id not set")
		}

		result, err := collector.Similar(ctx.Context, models.ApplicationSpec{
			AppID:   appIDFlag.Get(ctx),
			Lang:    langFlag.Get(ctx),
			Country: countryFlag.Get(ctx),
		})
		if err != nil {
			return fmt.Errorf("similar method: %w", err)
		}

		log.Ctx(ctx.Context).Debug().Msgf("similar apps founded: %d", len(result))
		display(ctx, result, appTable)

		return nil
	},
}

func appTable(apps []models.App) ([]string, [][]string) {
	var rows [][]string
	for _, app := range apps {
		rows = append(rows, []string{
			app.Title,
			app.URL,
			func() string {
				if app.Free {
					return "Free"
				}

				return fmt.Sprintf("%.2f %s", app.Price, app.Currency)
			}(),
			app.ScoreText,
			app.Developer,
		})
	}

	return []string{
		"Title",
		"URL",
		"Price",
		"Score",
		"Developer",
	}, rows
}