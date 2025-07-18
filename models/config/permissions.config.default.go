package config

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"
	authconfig "github.com/a-novel/service-authentication/models/config"
)

//go:embed permissions.config.yaml
var defaultPermissionsFile []byte

var PermissionsConfigDefault = config.MustUnmarshal[authconfig.Permissions](yaml.Unmarshal, defaultPermissionsFile)
