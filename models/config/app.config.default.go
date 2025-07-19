package config

import (
	"time"

	"github.com/samber/lo"

	"github.com/a-novel/golib/config"
	otelpresets "github.com/a-novel/golib/otel/presets"
	"github.com/a-novel/golib/postgres"
	"github.com/a-novel/golib/smtp"
)

const (
	SentryFlushTimeout = 2 * time.Second

	AppName = "service-story-schematics"

	APIPort                 = 8080
	APITimeoutRead          = 5 * time.Second
	APITimeoutReadHeader    = 3 * time.Second
	APITimeoutWrite         = 10 * time.Second
	APITimeoutIdle          = 30 * time.Second
	APITimeoutRequest       = 15 * time.Second
	APIMaxRequestSize       = 2 << 20 // 2 MiB
	APICorsAllowCredentials = false
	APICorsMaxAge           = 3600
)

var (
	APICorsAllowedOrigins = []string{"*"}
	APICorsAllowedHeaders = []string{"*"}
)

var isDebug = config.LoadEnv(
	lo.CoalesceOrEmpty(getEnv("SENTRY_DEBUG"), getEnv("DEBUG")), false, config.BoolParser,
)

var SMTPProd = smtp.ProdSender{
	Addr:     getEnv("SMTP_ADDR"),
	Name:     getEnv("SMTP_SENDER_NAME"),
	Email:    getEnv("SMTP_SENDER_EMAIL"),
	Password: getEnv("SMTP_SENDER_PASSWORD"),
	Domain:   getEnv("SMTP_SENDER_DOMAIN"),
}

var AppPresetDefault = App[*otelpresets.SentryOtelConfig, postgres.Config]{
	App: Main{
		Name: config.LoadEnv(getEnv("APP_NAME"), AppName, config.StringParser),
	},
	API: API{
		Port:           config.LoadEnv(getEnv("API_PORT"), APIPort, config.IntParser),
		MaxRequestSize: config.LoadEnv(getEnv("API_MAX_REQUEST_SIZE"), APIMaxRequestSize, config.Int64Parser),
		Timeouts: APITimeouts{
			Read: config.LoadEnv(getEnv("API_TIMEOUT_READ"), APITimeoutRead, config.DurationParser),
			ReadHeader: config.LoadEnv(
				getEnv("API_TIMEOUT_READ_HEADER"), APITimeoutReadHeader, config.DurationParser,
			),
			Write:   config.LoadEnv(getEnv("API_TIMEOUT_WRITE"), APITimeoutWrite, config.DurationParser),
			Idle:    config.LoadEnv(getEnv("API_TIMEOUT_IDLE"), APITimeoutIdle, config.DurationParser),
			Request: config.LoadEnv(getEnv("API_TIMEOUT_REQUEST"), APITimeoutRequest, config.DurationParser),
		},
		Cors: Cors{
			AllowedOrigins: config.LoadEnv(
				getEnv("API_CORS_ALLOWED_ORIGINS"), APICorsAllowedOrigins, config.SliceParser(config.StringParser),
			),
			AllowedHeaders: config.LoadEnv(
				getEnv("API_CORS_ALLOWED_HEADERS"), APICorsAllowedHeaders, config.SliceParser(config.StringParser),
			),
			AllowCredentials: config.LoadEnv(
				getEnv("API_CORS_ALLOW_CREDENTIALS"), APICorsAllowCredentials, config.BoolParser,
			),
			MaxAge: config.LoadEnv(getEnv("API_CORS_MAX_AGE"), APICorsMaxAge, config.IntParser),
		},
	},

	DependenciesConfig: Dependencies{
		JSONKeysURL: getEnv("JSON_KEYS_SERVICE_URL"),
	},
	PermissionsConfig: PermissionsConfigDefault,

	OpenAI: OpenAIPresetDefault,
	Otel: &otelpresets.SentryOtelConfig{
		DSN:          getEnv("SENTRY_DSN"),
		ServerName:   config.LoadEnv(getEnv("APP_NAME"), AppName, config.StringParser),
		Release:      getEnv("SENTRY_RELEASE"),
		Environment:  lo.CoalesceOrEmpty(getEnv("SENTRY_ENVIRONMENT"), getEnv("ENV")),
		FlushTimeout: config.LoadEnv(getEnv("SENTRY_FLUSH_TIMEOUT"), SentryFlushTimeout, config.DurationParser),
		Debug:        isDebug,
	},
	Postgres: PostgresPresetDefault,
}
