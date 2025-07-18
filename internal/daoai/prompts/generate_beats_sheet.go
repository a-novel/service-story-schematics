package prompts

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed generate_beats_sheet.en.yaml
var generateBeatsSheetEnFile []byte

//go:embed generate_beats_sheet.fr.yaml
var generateBeatsSheetFrFile []byte

type GenerateBeatsSheetsType struct {
	System string `yaml:"system"`
}

var GenerateBeatsSheetEN = config.MustUnmarshal[GenerateBeatsSheetsType](yaml.Unmarshal, generateBeatsSheetEnFile)

var GenerateBeatsSheetFR = config.MustUnmarshal[GenerateBeatsSheetsType](yaml.Unmarshal, generateBeatsSheetFrFile)

var GenerateBeatsSheet = map[models.Lang]GenerateBeatsSheetsType{
	models.LangEN: GenerateBeatsSheetEN,
	models.LangFR: GenerateBeatsSheetFR,
}
