package testdata

import (
	_ "embed"

	"github.com/a-novel/golib/config"
	"github.com/goccy/go-yaml"
)

//go:embed generate_loglines.en.yaml
var generateLoglinesEnFile []byte

type GenerateLoglinesTestCase struct {
	Count int    `yaml:"count"`
	Theme string `yaml:"theme"`
}

type GenerateLoglinesPromptsType struct {
	Cases      map[string]GenerateLoglinesTestCase `yaml:"cases"`
	CheckAgent struct {
		Themed string `yaml:"themed"`
		Random string `yaml:"random"`
	} `yaml:"checkAgent"`
}

var GenerateLoglinesPrompt = config.MustUnmarshal[GenerateLoglinesPromptsType](yaml.Unmarshal, generateLoglinesEnFile)
