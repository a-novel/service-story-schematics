package testdata

import (
	_ "embed"
	"github.com/a-novel-kit/configurator"
	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/models"
)

//go:embed generate_beats_sheet.en.yaml
var generateBeatsSheetEnFile []byte

//go:embed generate_beats_sheet.fr.yaml
var generateBeatsSheetFrFile []byte

type GenerateBeatsSheetTestCase struct {
	Logline string           `yaml:"logline"`
	Plan    models.StoryPlan `yaml:"plan"`
}

type GenerateBeatsSheetPromptsType struct {
	Cases      map[string]GenerateBeatsSheetTestCase `yaml:"cases"`
	CheckAgent string                                `yaml:"checkAgent"`
}

var GenerateBeatsSheetPromptEN = configurator.NewLoader[GenerateBeatsSheetPromptsType](config.Loader).MustLoad(
	configurator.NewConfig("", generateBeatsSheetEnFile),
)

var GenerateBeatsSheetPromptFR = configurator.NewLoader[GenerateBeatsSheetPromptsType](config.Loader).MustLoad(
	configurator.NewConfig("", generateBeatsSheetFrFile),
)

var GenerateBeatsSheetPrompts = map[models.Lang]GenerateBeatsSheetPromptsType{
	models.LangEN: GenerateBeatsSheetPromptEN,
	models.LangFR: GenerateBeatsSheetPromptFR,
}
