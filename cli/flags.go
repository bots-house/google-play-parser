package cli

import "github.com/urfave/cli/v2"

var appIDFlag = &cli.StringFlag{
	Name:    "app-id",
	Aliases: []string{"id"},
	Usage:   "set app id such 'com.best.app'",
}

var countryFlag = &cli.StringFlag{
	Name:    "country",
	Aliases: []string{"c"},
	Usage:   "set country in ISO format",
}

var langFlag = &cli.StringFlag{
	Name:    "lang",
	Aliases: []string{"l"},
	Usage:   "set lang",
}
