package testdata

import (
	_ "embed"
	"github.com/a-novel/golib/config"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/goccy/go-yaml"
)

//go:embed expand_logline.en.yaml
var expandLoglineEnFile []byte

//go:embed expand_logline.fr.yaml
var expandLoglineFrFile []byte

type ExpandLoglineTestCase struct {
	Logline string `yaml:"logline"`
}

type ExpandLoglinePromptsType struct {
	Cases      map[string]ExpandLoglineTestCase `yaml:"cases"`
	CheckAgent string                           `yaml:"checkAgent"`
}

var ExpandLoglinePromptEN = config.MustUnmarshal[ExpandLoglinePromptsType](yaml.Unmarshal, expandLoglineEnFile)

var ExpandLoglinePromptFR = config.MustUnmarshal[ExpandLoglinePromptsType](yaml.Unmarshal, expandLoglineFrFile)

var ExpandLoglinePrompts = map[models.Lang]ExpandLoglinePromptsType{
	models.LangEN: ExpandLoglinePromptEN,
	models.LangFR: ExpandLoglinePromptFR,
}
