package testdata

import (
	_ "embed"
	"github.com/a-novel/golib/config"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/goccy/go-yaml"
)

//go:embed regenerate_beats.en.yaml
var regenerateBeatsEnFile []byte

//go:embed regenerate_beats.fr.yaml
var regenerateBeatsFrFile []byte

type RegenerateBeatsTestCase struct {
	Logline        string           `yaml:"logline"`
	Plan           models.StoryPlan `yaml:"plan"`
	Beats          []models.Beat    `yaml:"beats"`
	RegenerateKeys []string         `yaml:"regenerateKeys"`
}

type RegenerateBeatsPromptsType struct {
	Cases      map[string]RegenerateBeatsTestCase `yaml:"cases"`
	CheckAgent string                             `yaml:"checkAgent"`
}

var RegenerateBeatsPromptEN = config.MustUnmarshal[RegenerateBeatsPromptsType](yaml.Unmarshal, regenerateBeatsEnFile)

var RegenerateBeatsPromptFR = config.MustUnmarshal[RegenerateBeatsPromptsType](yaml.Unmarshal, regenerateBeatsFrFile)

var RegenerateBeatsPrompts = map[models.Lang]RegenerateBeatsPromptsType{
	models.LangEN: RegenerateBeatsPromptEN,
	models.LangFR: RegenerateBeatsPromptFR,
}
