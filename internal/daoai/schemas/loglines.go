package schemas

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed loglines.en.yaml
var loglinesEnFile []byte

//go:embed loglines.fr.yaml
var loglinesFrFile []byte

var LoglinesEN = config.MustUnmarshal[Schema](yaml.Unmarshal, loglinesEnFile)

var LoglinesFR = config.MustUnmarshal[Schema](yaml.Unmarshal, loglinesFrFile)

var Loglines = map[models.Lang]Schema{
	models.LangEN: LoglinesEN,
	models.LangFR: LoglinesFR,
}
