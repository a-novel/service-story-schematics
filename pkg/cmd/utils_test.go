package cmdpkg_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	otelpresets "github.com/a-novel/golib/otel/presets"
	"github.com/a-novel/golib/postgres"
	authmodels "github.com/a-novel/service-authentication/models"
	jkmodels "github.com/a-novel/service-json-keys/models"
	jkpkg "github.com/a-novel/service-json-keys/pkg"

	"github.com/a-novel/service-story-schematics/models/config"
)

type TestConfig = config.App[*otelpresets.LocalOtelConfig, postgres.Config]

func getAccessToken(t *testing.T, appConfig TestConfig, claims authmodels.AccessTokenClaims) string {
	t.Helper()

	jsonKeysClient, err := jkpkg.NewAPIClient(context.Background(), appConfig.DependenciesConfig.JSONKeysURL)
	require.NoError(t, err)

	signer := jkpkg.NewClaimsSigner(jsonKeysClient)

	accessToken, err := signer.SignClaims(
		context.Background(),
		jkmodels.KeyUsageAuth,
		claims,
	)
	require.NoError(t, err)

	return accessToken
}
