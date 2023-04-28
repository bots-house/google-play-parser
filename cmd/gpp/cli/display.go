package cli

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
)

func display[T any](c *cli.Context, data T) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if _, err := c.App.Writer.Write(body); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
