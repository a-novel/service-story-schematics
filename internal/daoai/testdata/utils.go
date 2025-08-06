package testdata

import (
	_ "embed"

	"github.com/a-novel/golib/config"
	"github.com/goccy/go-yaml"
)

//go:embed utils.en.yaml
var utilsEnFile []byte

type UtilsPromptsType struct {
	CheckAgent struct {
		System string `yaml:"system"`
		Expect string `yaml:"expect"`
	} `yaml:"checkAgent"`
}

var UtilsPrompt = config.MustUnmarshal[UtilsPromptsType](yaml.Unmarshal, utilsEnFile)
