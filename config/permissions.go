package config

import (
	_ "embed"

	"github.com/a-novel/authentication/models"

	"github.com/a-novel-kit/configurator"
)

//go:embed permissions.yaml
var permissionsFile []byte

var Permissions = configurator.NewLoader[models.PermissionsConfig](Loader).
	MustLoad(configurator.NewConfig("", permissionsFile))
