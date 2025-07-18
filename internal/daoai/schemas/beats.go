package schemas

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed beats.en.yaml
var beatsEnFile []byte

//go:embed beats.fr.yaml
var beatsFrFile []byte

var BeatsEN = config.MustUnmarshal[Schema](yaml.Unmarshal, beatsEnFile)

var BeatsFR = config.MustUnmarshal[Schema](yaml.Unmarshal, beatsFrFile)

var Beats = map[models.Lang]Schema{
	models.LangEN: BeatsEN,
	models.LangFR: BeatsFR,
}
