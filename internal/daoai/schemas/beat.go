package schemas

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed beat.en.yaml
var beatEnFile []byte

//go:embed beat.fr.yaml
var beatFrFile []byte

var BeatEN = config.MustUnmarshal[Schema](yaml.Unmarshal, beatEnFile)

var BeatFR = config.MustUnmarshal[Schema](yaml.Unmarshal, beatFrFile)

var Beat = map[models.Lang]Schema{
	models.LangEN: BeatEN,
	models.LangFR: BeatFR,
}
