package prompts

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"
)

//go:embed regenerate_beats.en.yaml
var regenerateBeatsEnFile []byte

type RegenerateBeatsType struct {
	System string `yaml:"system"`
	Input1 string `yaml:"input1"`
	Input2 string `yaml:"input2"`
}

var RegenerateBeats = config.MustUnmarshal[RegenerateBeatsType](yaml.Unmarshal, regenerateBeatsEnFile)
