package schemas

import (
	_ "embed"

	"github.com/a-novel-kit/configurator"

	"github.com/a-novel/service-story-schematics/config"
)

//go:embed en.yaml
var en []byte

//go:embed fr.yaml
var fr []byte

type Schema struct {
	Description string `yaml:"description"`
	Schema      any    `yaml:"schema"`
}

type Schemas struct {
	Beat     Schema `yaml:"beat"`
	Beats    Schema `yaml:"beats"`
	Logline  Schema `yaml:"logline"`
	Loglines Schema `yaml:"loglines"`
}

type TranslatedSchemas struct {
	En Schemas `yaml:"en"`
	Fr Schemas `yaml:"fr"`
}

var Config = TranslatedSchemas{
	En: configurator.NewLoader[Schemas](config.Loader).MustLoad(configurator.NewConfig("", en)),
	Fr: configurator.NewLoader[Schemas](config.Loader).MustLoad(configurator.NewConfig("", fr)),
}
