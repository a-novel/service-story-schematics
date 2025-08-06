package testdata

import (
	_ "embed"

	"github.com/a-novel/golib/config"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/goccy/go-yaml"
)

//go:embed generate_beats_sheet.en.yaml
var generateBeatsSheetEnFile []byte

type GenerateBeatsSheetTestCase struct {
	Logline string           `yaml:"logline"`
	Plan    models.StoryPlan `yaml:"plan"`
}

type GenerateBeatsSheetPromptsType struct {
	Cases      map[string]GenerateBeatsSheetTestCase `yaml:"cases"`
	CheckAgent string                                `yaml:"checkAgent"`
}

var GenerateBeatsSheetPrompt = config.MustUnmarshal[GenerateBeatsSheetPromptsType](
	yaml.Unmarshal, generateBeatsSheetEnFile,
)
