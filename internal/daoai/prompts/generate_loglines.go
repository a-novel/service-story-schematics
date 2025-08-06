package prompts

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"
)

//go:embed generate_loglines.en.yaml
var generateLoglinesEnFile []byte

type GenerateLoglinessType struct {
	System struct {
		Themed string `yaml:"themed"`
		Random string `yaml:"random"`
	} `yaml:"system"`
}

var GenerateLoglines = config.MustUnmarshal[GenerateLoglinessType](yaml.Unmarshal, generateLoglinesEnFile)
