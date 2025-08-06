package schemas

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"
)

//go:embed beats.en.yaml
var beatsEnFile []byte

var Beats = config.MustUnmarshal[Schema](yaml.Unmarshal, beatsEnFile)
