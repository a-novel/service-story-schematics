package prompts

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed langs.yaml
var langsFile []byte

type LangsType map[models.Lang]string

var Langs = config.MustUnmarshal[LangsType](yaml.Unmarshal, langsFile)
