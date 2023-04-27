package cli

import (
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

func display[T any](c *cli.Context, data T, tableFunc func(T) ([]string, [][]string)) {
	header, rows := tableFunc(data)

	table := tablewriter.NewWriter(c.App.Writer)

	table.SetHeader(header)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetHeaderLine(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.AppendBulk(rows)
	table.Render()
}
