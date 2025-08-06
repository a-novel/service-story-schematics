package testdata

import (
	_ "embed"

	"github.com/a-novel/golib/config"
	"github.com/goccy/go-yaml"
)

//go:embed expand_logline.en.yaml
var expandLoglineEnFile []byte

type ExpandLoglineTestCase struct {
	Logline string `yaml:"logline"`
}

type ExpandLoglinePromptsType struct {
	Cases      map[string]ExpandLoglineTestCase `yaml:"cases"`
	CheckAgent string                           `yaml:"checkAgent"`
}

var ExpandLoglinePrompt = config.MustUnmarshal[ExpandLoglinePromptsType](yaml.Unmarshal, expandLoglineEnFile)
