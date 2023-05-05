package gpp

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"

	"github.com/bots-house/google-play-parser/models"
)

func checkApp(app *App) error {
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
		if !assert.NoError(t, err) {
			return
		}

		for _, app := range apps {
			assert.NoError(t, checkApp(&app))
		}
	})

	t.Run("App", func(t *testing.T) {
		app, err := collector.App(context.Background(), ApplicationSpec{
			AppID: "com.mojang.minecraftpe",
		})
		if !assert.NoError(t, err) {
			return
		}

		assert.NoError(t, checkApp(&app))
	})

	t.Run("List", func(t *testing.T) {
		count := rand.Intn(100)

		tests := []struct {
			name      string
			spec      ListSpec
			assertion func(apps []App, err error, wantErr bool)
			wantErr   bool
		}{
			{
				name: "Simple",
				spec: ListSpec{Count: count},
				assertion: func(apps []App, err error, wantErr bool) {
					if err != nil && !wantErr {
						assert.NoError(t, err)
						return
					}

					assert.Equal(t, count, len(apps))

					for _, app := range apps {
						assert.NoError(t, checkApp(&app))
					}
				},
			},

			{
				name: "WithoutCount",
				spec: ListSpec{},
				assertion: func(apps []App, err error, wantErr bool) {
					if err != nil && !wantErr {
						assert.NoError(t, err)
						return
					}

					assert.Equal(t, models.GetDefaultListCount(), len(apps))

					for _, app := range apps {
						assert.NoError(t, checkApp(&app))
					}
				},
			},

			{
				name: "InvalidSpec",
				spec: ListSpec{Count: -1},
				assertion: func(apps []App, err error, wantErr bool) {
					assert.Error(t, err)
					assert.True(t, wantErr)
				},
				wantErr: true,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				apps, err := collector.List(context.Background(), test.spec)
				test.assertion(apps, err, test.wantErr)
			})
		}
	})

	t.Run("Developer", func(t *testing.T) {
		tests := []struct {
			name string
			spec DeveloperSpec
		}{
			{
				name: "Developer name",
				spec: DeveloperSpec{
					DevID: "Jam City, Inc.",
				},
			},

			{
				name: "Developer id",
				spec: DeveloperSpec{
					DevID: "5700313618786177705",
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				apps, err := collector.Developer(context.Background(), test.spec)
				if !assert.NoError(t, err) {
					return
				}

				for _, app := range apps {
					assert.NoError(t, checkApp(&app))
				}
			})
		}
	})

	t.Run("Search", func(t *testing.T) {
		apps, err := collector.Search(context.Background(), SearchSpec{Query: "netflix", Count: 1})
		if !assert.NoError(t, err) {
			return
		}

		for _, app := range apps {
			assert.NoError(t, checkApp(&app))
		}

		// Check for main app
		assert.Equal(t, "Netflix", apps[0].Title)
		assert.Equal(t, "com.netflix.mediaclient", apps[0].AppID)
	})
}
