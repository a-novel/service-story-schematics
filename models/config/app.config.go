package config

import (
	"time"

	"github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/postgres"
	authconfig "github.com/a-novel/service-authentication/models/config"
)

type Main struct {
	Name string `json:"name" yaml:"name"`
}

type Dependencies struct {
	JSONKeysURL string `json:"jsonKeysURL" yaml:"jsonKeysURL"`
}

type APITimeouts struct {
	Read       time.Duration `json:"read"       yaml:"read"`
	ReadHeader time.Duration `json:"readHeader" yaml:"readHeader"`
	Write      time.Duration `json:"write"      yaml:"write"`
	Idle       time.Duration `json:"idle"       yaml:"idle"`
	Request    time.Duration `json:"request"    yaml:"request"`
}

type Cors struct {
	AllowedOrigins   []string `json:"allowedOrigins"   yaml:"allowedOrigins"`
	AllowedHeaders   []string `json:"allowedHeaders"   yaml:"allowedHeaders"`
	AllowCredentials bool     `json:"allowCredentials" yaml:"allowCredentials"`
	MaxAge           int      `json:"maxAge"           yaml:"maxAge"`
}

type API struct {
	Port           int         `json:"port"           yaml:"port"`
	Timeouts       APITimeouts `json:"timeouts"       yaml:"timeouts"`
	MaxRequestSize int64       `json:"maxRequestSize" yaml:"maxRequestSize"`
	Cors           Cors        `json:"cors"           yaml:"cors"`
}

type App[Otel otel.Config, Pg postgres.Config] struct {
	App Main `json:"app" yaml:"app"`
	API API  `json:"api" yaml:"api"`

	DependenciesConfig Dependencies           `json:"dependencies" yaml:"dependencies"`
	PermissionsConfig  authconfig.Permissions `json:"permissions"  yaml:"permissions"`

	OpenAI   OpenAI `json:"openai"   yaml:"openai"`
	Otel     Otel   `json:"otel"     yaml:"otel"`
	Postgres Pg     `json:"postgres" yaml:"postgres"`
}
