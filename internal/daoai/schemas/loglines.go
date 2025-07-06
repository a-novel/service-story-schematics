package schemas

import (
	_ "embed"

	"github.com/a-novel-kit/configurator"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/models"
)

//go:embed loglines.en.yaml
var loglinesEnFile []byte

//go:embed loglines.fr.yaml
var loglinesFrFile []byte

var LoglinesEN = configurator.NewLoader[Schema](config.Loader).MustLoad(
	configurator.NewConfig("", loglinesEnFile),
)

var LoglinesFR = configurator.NewLoader[Schema](config.Loader).MustLoad(
	configurator.NewConfig("", loglinesFrFile),
)

var Loglines = map[models.Lang]Schema{
	models.LangEN: LoglinesEN,
	models.LangFR: LoglinesFR,
}
