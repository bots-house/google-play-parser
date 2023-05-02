package gpp

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

func checkApp(app App) error {
	errs := make([]error, 0, 8)

	if app.Title == "" {
		errs = append(errs, fmt.Errorf("title not parsed"))
	}

	if app.AppID == "" {
		errs = append(errs, fmt.Errorf("appID not parsed"))
	}

	if app.URL == "" {
		errs = append(errs, fmt.Errorf("url not parsed"))
	}

	if app.Icon == "" {
		errs = append(errs, fmt.Errorf("icon not parsed"))
	}

	if app.Developer == "" {
		errs = append(errs, fmt.Errorf("developer not parsed"))
	}

	if app.Price == 0 && !app.Free {
		errs = append(errs, fmt.Errorf("price not parsed"))
	}

	if app.Summary == "" {
		errs = append(errs, fmt.Errorf("summary not parsed"))
	}

	if app.Score == 0 && app.ScoreText == "" {
		errs = append(errs, fmt.Errorf("score not parsed"))
	}

	return multierr.Combine(errs...)
}

func Test_Scraper(t *testing.T) {
	collector := New()

	t.Run("Similar", func(t *testing.T) {
		apps, err := collector.Similar(
			context.Background(),
			ApplicationSpec{AppID: "com.mojang.minecraftpe"},
		)
		if err != nil {
			t.Fatal(err)
		}

		for _, app := range apps {
			assert.NoError(t, checkApp(app))
		}
	})

	t.Run("App", func(t *testing.T) {
		app, err := collector.App(context.Background(), ApplicationSpec{
			AppID: "com.mojang.minecraftpe",
		})
		if err != nil {
			t.Fatal(err)
		}

		assert.NoError(t, checkApp(app))
	})
}
