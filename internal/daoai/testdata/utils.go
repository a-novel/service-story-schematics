package testdata

import (
	_ "embed"
	"github.com/a-novel/golib/config"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/goccy/go-yaml"
)

//go:embed utils.en.yaml
var utilsEnFile []byte

//go:embed utils.fr.yaml
var utilsFrFile []byte

type UtilsPromptsType struct {
	CheckAgent struct {
		System string `yaml:"system"`
		Expect string `yaml:"expect"`
	} `yaml:"checkAgent"`
}

var UtilsPromptEN = config.MustUnmarshal[UtilsPromptsType](yaml.Unmarshal, utilsEnFile)

var UtilsPromptFR = config.MustUnmarshal[UtilsPromptsType](yaml.Unmarshal, utilsFrFile)

var UtilsPrompts = map[models.Lang]UtilsPromptsType{
	models.LangEN: UtilsPromptEN,
	models.LangFR: UtilsPromptFR,
}
