package prompts

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"
)

//go:embed expand_logline.en.yaml
var expandLoglineEnFile []byte

type ExpandLoglinesType struct {
	System string `yaml:"system"`
}

var ExpandLogline = config.MustUnmarshal[ExpandLoglinesType](yaml.Unmarshal, expandLoglineEnFile)
