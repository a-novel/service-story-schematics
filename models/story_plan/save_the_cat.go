package storyplanmodel

import (
	_ "embed"

	"github.com/goccy/go-yaml"
	"github.com/samber/lo"

	"github.com/a-novel/golib/config"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed save_the_cat.en.yaml
var saveTheCatEN []byte

//go:embed save_the_cat.fr.yaml
var saveTheCatFR []byte

var SaveTheCat = map[models.Lang]*Plan{
	models.LangEN: lo.ToPtr(config.MustUnmarshal[Plan](yaml.Unmarshal, saveTheCatEN)),
	models.LangFR: lo.ToPtr(config.MustUnmarshal[Plan](yaml.Unmarshal, saveTheCatFR)),
}
