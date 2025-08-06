package schemas

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"
)

//go:embed loglines.en.yaml
var loglinesEnFile []byte

var Loglines = config.MustUnmarshal[Schema](yaml.Unmarshal, loglinesEnFile)
