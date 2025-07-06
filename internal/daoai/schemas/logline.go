package schemas

import (
	_ "embed"

	"github.com/a-novel-kit/configurator"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/models"
)

//go:embed logline.en.yaml
var loglineEnFile []byte

//go:embed logline.fr.yaml
var loglineFrFile []byte

var LoglineEN = configurator.NewLoader[Schema](config.Loader).MustLoad(
	configurator.NewConfig("", loglineEnFile),
)

var LoglineFR = configurator.NewLoader[Schema](config.Loader).MustLoad(
	configurator.NewConfig("", loglineFrFile),
)

var Logline = map[models.Lang]Schema{
	models.LangEN: LoglineEN,
	models.LangFR: LoglineFR,
}
