package schemas

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"
)

//go:embed beat.en.yaml
var beatEnFile []byte

var Beat = config.MustUnmarshal[Schema](yaml.Unmarshal, beatEnFile)
