package prompts

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed regenerate_beats.en.yaml
var regenerateBeatsEnFile []byte

//go:embed regenerate_beats.fr.yaml
var regenerateBeatsFrFile []byte

type RegenerateBeatsType struct {
	System string `yaml:"system"`
	Input1 string `yaml:"input1"`
	Input2 string `yaml:"input2"`
}

var RegenerateBeatsEN = config.MustUnmarshal[RegenerateBeatsType](yaml.Unmarshal, regenerateBeatsEnFile)

var RegenerateBeatsFR = config.MustUnmarshal[RegenerateBeatsType](yaml.Unmarshal, regenerateBeatsFrFile)

var RegenerateBeats = map[models.Lang]RegenerateBeatsType{
	models.LangEN: RegenerateBeatsEN,
	models.LangFR: RegenerateBeatsFR,
}
