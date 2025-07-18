package config

import (
	"os"

	"github.com/samber/lo"

	"github.com/a-novel/golib/config"
	otelpresets "github.com/a-novel/golib/otel/presets"
	"github.com/a-novel/golib/postgres"
)

func AppPresetTest(port int) App[*otelpresets.SentryOtelConfig, postgres.Config] {
	return App[*otelpresets.SentryOtelConfig, postgres.Config]{
		App: Main{
			Name: config.LoadEnv(os.Getenv("APP_NAME"), AppName, config.StringParser),
		},
		API: API{
			Port:           port,
			MaxRequestSize: config.LoadEnv(os.Getenv("API_MAX_REQUEST_SIZE"), APIMaxRequestSize, config.Int64Parser),
			Timeouts: APITimeouts{
				Read: config.LoadEnv(os.Getenv("API_TIMEOUT_READ"), APITimeoutRead, config.DurationParser),
				ReadHeader: config.LoadEnv(
					os.Getenv("API_TIMEOUT_READ_HEADER"), APITimeoutReadHeader, config.DurationParser,
				),
				Write:   config.LoadEnv(os.Getenv("API_TIMEOUT_WRITE"), APITimeoutWrite, config.DurationParser),
				Idle:    config.LoadEnv(os.Getenv("API_TIMEOUT_IDLE"), APITimeoutIdle, config.DurationParser),
				Request: config.LoadEnv(os.Getenv("API_TIMEOUT_REQUEST"), APITimeoutRequest, config.DurationParser),
			},
			Cors: Cors{
				AllowedOrigins: config.LoadEnv(
					os.Getenv("API_CORS_ALLOWED_ORIGINS"), APICorsAllowedOrigins, config.SliceParser(config.StringParser),
				),
				AllowedHeaders: config.LoadEnv(
					os.Getenv("API_CORS_ALLOWED_HEADERS"), APICorsAllowedHeaders, config.SliceParser(config.StringParser),
				),
				AllowCredentials: config.LoadEnv(
					os.Getenv("API_CORS_ALLOW_CREDENTIALS"), APICorsAllowCredentials, config.BoolParser,
				),
				MaxAge: config.LoadEnv(os.Getenv("API_CORS_MAX_AGE"), APICorsMaxAge, config.IntParser),
			},
		},

		DependenciesConfig: Dependencies{
			JSONKeysURL: os.Getenv("JSON_KEYS_SERVICE_TEST_URL"),
		},
		PermissionsConfig: PermissionsConfigDefault,

		OpenAI: OpenAIPresetDefault,
		Otel: &otelpresets.SentryOtelConfig{
			DSN:          os.Getenv("SENTRY_DSN"),
			ServerName:   config.LoadEnv(os.Getenv("APP_NAME"), AppName, config.StringParser),
			Release:      os.Getenv("SENTRY_RELEASE"),
			Environment:  lo.CoalesceOrEmpty(os.Getenv("SENTRY_ENVIRONMENT"), os.Getenv("ENV")),
			FlushTimeout: config.LoadEnv(os.Getenv("SENTRY_FLUSH_TIMEOUT"), SentryFlushTimeout, config.DurationParser),
			Debug:        isDebug,
		},
		Postgres: PostgresPresetTest,
	}
}
