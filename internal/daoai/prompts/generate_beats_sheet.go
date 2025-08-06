package prompts

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"
)

//go:embed generate_beats_sheet.en.yaml
var generateBeatsSheetEnFile []byte

type GenerateBeatsSheetsType struct {
	System string `yaml:"system"`
}

var GenerateBeatsSheet = config.MustUnmarshal[GenerateBeatsSheetsType](yaml.Unmarshal, generateBeatsSheetEnFile)
