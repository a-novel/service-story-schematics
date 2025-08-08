package testdata

import (
	_ "embed"

	"github.com/a-novel/golib/config"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/goccy/go-yaml"
)

//go:embed expand_beat.en.yaml
var expandBeatEnFile []byte

type ExpandBeatTestCase struct {
	Logline   string        `yaml:"logline"`
	Beats     []models.Beat `yaml:"beats"`
	TargetKey string        `yaml:"targetKey"`
}

type ExpandBeatPromptsType struct {
	Cases      map[string]ExpandBeatTestCase `yaml:"cases"`
	CheckAgent string                        `yaml:"checkAgent"`
}

var ExpandBeatPrompt = config.MustUnmarshal[ExpandBeatPromptsType](yaml.Unmarshal, expandBeatEnFile)
