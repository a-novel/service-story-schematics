package schemas

import (
	_ "embed"

	"github.com/a-novel-kit/configurator"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/models"
)

//go:embed beat.en.yaml
var beatEnFile []byte

//go:embed beat.fr.yaml
var beatFrFile []byte

var BeatEN = configurator.NewLoader[Schema](config.Loader).MustLoad(
	configurator.NewConfig("", beatEnFile),
)

var BeatFR = configurator.NewLoader[Schema](config.Loader).MustLoad(
	configurator.NewConfig("", beatFrFile),
)

var Beat = map[models.Lang]Schema{
	models.LangEN: BeatEN,
	models.LangFR: BeatFR,
}
