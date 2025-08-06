package schemas

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"
)

//go:embed logline.en.yaml
var loglineEnFile []byte

var Logline = config.MustUnmarshal[Schema](yaml.Unmarshal, loglineEnFile)
