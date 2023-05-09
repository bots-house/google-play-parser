package cli

import "github.com/urfave/cli/v2"

var appIDFlag = &cli.StringFlag{
	Name:    "app-id",
	Aliases: []string{"id"},
	Usage:   "set app id such 'com.best.app'",
}

var devIDFlag = &cli.StringFlag{
	Name:  "dev-id",
	Usage: "set dev id such",
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

var categoryFlag = &cli.StringFlag{
	Name:  "category",
	Usage: "set category",
}

var collectionFlag = &cli.StringFlag{
	Name:  "collection",
	Usage: "set collection",
}

var countFlag = &cli.IntFlag{
	Name:  "count",
	Usage: "set count for returned entries",
}

var ageFlag = &cli.StringFlag{
	Name:  "age",
	Usage: "set collection",
}

var fullFlag = &cli.BoolFlag{
	Name:  "full",
	Usage: "if true parse full detail",
}

var queryFlag = &cli.StringFlag{
	Name:    "query",
	Aliases: []string{"q"},
}

var priceFlag = &cli.StringFlag{
	Name:    "price",
	Aliases: []string{"p"},
}

var sortFlag = &cli.StringFlag{
	Name:  "sort",
	Usage: "Possible values: [helpfulness, newest, rating]",
}
