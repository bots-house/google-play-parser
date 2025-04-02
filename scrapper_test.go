package gpp

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"

	"github.com/bots-house/google-play-parser/internal/shared"
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
					Full:  true,
				},
			},

			{
				name: "Developer id",
				spec: DeveloperSpec{
					DevID: "5700313618786177705",
					Full:  true,
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

		assert.Len(t, apps, 1)

		for _, app := range apps {
			assert.NoError(t, checkApp(&app))
		}

		// Check for main app
		assert.Equal(t, "Netflix", apps[0].Title)
		assert.Equal(t, "com.netflix.mediaclient", apps[0].AppID)
	})

	t.Run("DataSafety", func(t *testing.T) {
		result, err := collector.DataSafety(context.Background(), ApplicationSpec{
			AppID: "com.sgn.pandapop.gp",
		})
		if !assert.NoError(t, err) {
			return
		}

		assert.NotEmpty(t, result.SharedData)
		assert.NotEmpty(t, result.CollectedData)
		assert.NotEmpty(t, result.SecurityPractice)

		_, err = url.Parse(result.PrivacyPolicyURL)
		assert.NoError(t, err)
	})

	t.Run("Permissions", func(t *testing.T) {
		perms, err := collector.Permissions(context.Background(), ApplicationSpec{
			AppID: "com.sgn.pandapop.gp",
			Full:  true,
		})
		if !assert.NoError(t, err) {
			return
		}

		assert.Len(t, perms, 14)
	})

	t.Run("Suggest", func(t *testing.T) {
		result, err := collector.Suggest(context.Background(), SearchSpec{Query: "p"})

		if !assert.NoError(t, err) {
			return
		}

		assert.Len(t, result, 5)

		for _, ok := range shared.Map(result, func(val string) bool {
			return strings.HasPrefix(val, "p")
		}) {
			assert.True(t, ok)
		}
	})

	t.Run("Reviews", func(t *testing.T) {
		result, err := collector.Reviews(context.Background(), ReviewsSpec{AppID: "com.sgn.pandapop.gp", Count: 1})
		if !assert.NoError(t, err) {
			return
		}

		assert.Len(t, result, 1)

		for _, review := range result {
			assert.NotEmpty(t, review.ID)
			assert.NotEmpty(t, review.URL)
		}
	})
}

func Test_InAppPurchases(t *testing.T) {
	tests := []struct {
		appID         string
		inAppPurchase bool
	}{
		{
			appID:         "com.mojang.minecraftpe",
			inAppPurchase: true,
		},

		{
			appID:         "com.miniclip.plagueinc",
			inAppPurchase: true,
		},

		{
			appID: "com.einnovation.temu",
		},

		{
			appID: "com.whatsapp",
		},
	}

	collector := New()

	for _, test := range tests {
		app, err := collector.App(context.Background(), ApplicationSpec{
			AppID: test.appID,
		})
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, test.inAppPurchase, app.InAppPurchase)
	}
}

func TestMissingDeveloperNames(t *testing.T) {
	c := New()

	tests := []struct {
		id string
	}{
		{
			id: "com.particlenews.newsbreak",
		},

		{
			id: "com.xphotokit.chatgptassist",
		},

		{
			id: "com.newleaf.app.android.victor",
		},
	}

	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			app, err := c.App(context.Background(), ApplicationSpec{AppID: test.id})
			if err != nil {
				t.Error(err)
				return
			}

			if _, err := strconv.ParseInt(app.Developer, 10, strconv.IntSize); err == nil {
				t.Error("developer name is id")
			}
		})
	}
}
