package schemas

import (
	_ "embed"

	"github.com/a-novel-kit/configurator"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/models"
)

//go:embed beats.en.yaml
var beatsEnFile []byte

//go:embed beats.fr.yaml
var beatsFrFile []byte

var BeatsEN = configurator.NewLoader[Schema](config.Loader).MustLoad(
	configurator.NewConfig("", beatsEnFile),
)

var BeatsFR = configurator.NewLoader[Schema](config.Loader).MustLoad(
	configurator.NewConfig("", beatsFrFile),
)

var Beats = map[models.Lang]Schema{
	models.LangEN: BeatsEN,
	models.LangFR: BeatsFR,
}
