package cli

import "github.com/urfave/cli/v2"

func New() *cli.App {
	return &cli.App{
		Name: "parser",
		Commands: []*cli.Command{
			appCMD,
			similarCMD,
			listCMD,
			developerCMD,
			searchCMD,
			dataSafetyCMD,
			permissionsCMD,
			suggestCMD,
			reviewsCMD,
		},
	}
}
