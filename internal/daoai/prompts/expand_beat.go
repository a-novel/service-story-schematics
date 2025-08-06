package prompts

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"
)

//go:embed expand_beat.en.yaml
var expandBeatEnFile []byte

type ExpandBeatsType struct {
	System string `yaml:"system"`
	Input1 string `yaml:"input1"`
	Input2 string `yaml:"input2"`
}

var ExpandBeat = config.MustUnmarshal[ExpandBeatsType](yaml.Unmarshal, expandBeatEnFile)
