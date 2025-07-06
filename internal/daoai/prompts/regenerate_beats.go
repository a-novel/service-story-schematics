package prompts

import (
	_ "embed"

	"github.com/a-novel-kit/configurator"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/models"
)

//go:embed regenerate_beats.en.yaml
var regenerateBeatsEnFile []byte

//go:embed regenerate_beats.fr.yaml
var regenerateBeatsFrFile []byte

type RegenerateBeatssType struct {
	System string `yaml:"system"`
	Input1 string `yaml:"input1"`
	Input2 string `yaml:"input2"`
}

var RegenerateBeatsEN = configurator.NewLoader[RegenerateBeatssType](config.Loader).MustLoad(
	configurator.NewConfig("", regenerateBeatsEnFile),
)

var RegenerateBeatsFR = configurator.NewLoader[RegenerateBeatssType](config.Loader).MustLoad(
	configurator.NewConfig("", regenerateBeatsFrFile),
)

var RegenerateBeats = map[models.Lang]RegenerateBeatssType{
	models.LangEN: RegenerateBeatsEN,
	models.LangFR: RegenerateBeatsFR,
}
