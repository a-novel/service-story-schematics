package prompts

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed expand_logline.en.yaml
var expandLoglineEnFile []byte

//go:embed expand_logline.fr.yaml
var expandLoglineFrFile []byte

type ExpandLoglinesType struct {
	System string `yaml:"system"`
}

var ExpandLoglineEN = config.MustUnmarshal[ExpandLoglinesType](yaml.Unmarshal, expandLoglineEnFile)

var ExpandLoglineFR = config.MustUnmarshal[ExpandLoglinesType](yaml.Unmarshal, expandLoglineFrFile)

var ExpandLogline = map[models.Lang]ExpandLoglinesType{
	models.LangEN: ExpandLoglineEN,
	models.LangFR: ExpandLoglineFR,
}
