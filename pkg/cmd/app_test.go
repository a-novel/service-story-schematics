package cmdpkg_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/ogen"
	"github.com/a-novel/golib/postgres"

	"github.com/a-novel/service-story-schematics/migrations"
	"github.com/a-novel/service-story-schematics/models/config"
	"github.com/a-novel/service-story-schematics/pkg"
	cmdpkg "github.com/a-novel/service-story-schematics/pkg/cmd"
)

type AppTestSuite func(ctx context.Context, t *testing.T, config TestConfig)

func TestApp(t *testing.T) {
	testSuites := map[string]AppTestSuite{
		"Ping":                  testAppPing,
		"BeatsSheetsPlayground": testAppBeatsSheetsPlayground,
		"LoglinesPlayground":    testAppLoglinesPlayground,
		"StoryPlansCRUD":        testAppStoryPlansCRUD,
	}

	for testName, testSuite := range testSuites {
		t.Run(testName, func(t *testing.T) {
			postgres.RunIsolatedTransactionalTest(
				t, config.PostgresPresetTest, migrations.Migrations, func(ctx context.Context, t *testing.T) {
					t.Helper()

					port, err := ogen.GetRandomPort()
					require.NoError(t, err)

					appConfig := config.AppPresetTest(port)

					go func() {
						assert.NoError(t, cmdpkg.App(ctx, appConfig))
					}()

					security := pkg.NewBearerSource()
					client, err := pkg.NewAPIClient(
						ctx, fmt.Sprintf("http://localhost:%v/v1", appConfig.API.Port), security,
					)
					require.NoError(t, err)

					require.Eventually(t, func() bool {
						_, err = client.Ping(t.Context())

						return assert.NoError(t, err)
					}, 10*time.Second, 100*time.Millisecond)

					testSuite(ctx, t, appConfig)
				},
			)
		})
	}
}
