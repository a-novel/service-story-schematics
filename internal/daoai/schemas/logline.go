package schemas

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed logline.en.yaml
var loglineEnFile []byte

//go:embed logline.fr.yaml
var loglineFrFile []byte

var LoglineEN = config.MustUnmarshal[Schema](yaml.Unmarshal, loglineEnFile)

var LoglineFR = config.MustUnmarshal[Schema](yaml.Unmarshal, loglineFrFile)

var Logline = map[models.Lang]Schema{
	models.LangEN: LoglineEN,
	models.LangFR: LoglineFR,
}
