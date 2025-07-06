package api_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/api"
	"github.com/a-novel/service-story-schematics/internal/api/codegen"
)

func TestPing(t *testing.T) {
	t.Parallel()

	res, err := new(api.API).Ping(t.Context())
	require.NoError(t, err)
	require.Equal(t, &codegen.PingOK{Data: strings.NewReader("pong")}, res)
}
