package prompts

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

type GenerateBeatsSheetsType struct {
	System string `yaml:"system"`
}

var GenerateBeatsSheetEN = configurator.NewLoader[GenerateBeatsSheetsType](config.Loader).MustLoad(
	configurator.NewConfig("", generateBeatsSheetEnFile),
)

var GenerateBeatsSheetFR = configurator.NewLoader[GenerateBeatsSheetsType](config.Loader).MustLoad(
	configurator.NewConfig("", generateBeatsSheetFrFile),
)

var GenerateBeatsSheet = map[models.Lang]GenerateBeatsSheetsType{
	models.LangEN: GenerateBeatsSheetEN,
	models.LangFR: GenerateBeatsSheetFR,
}
