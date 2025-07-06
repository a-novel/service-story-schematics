package main

import (
	"context"

	authModels "github.com/a-novel/service-authentication/models"
	jkModels "github.com/a-novel/service-json-keys/models"
	jkPkg "github.com/a-novel/service-json-keys/pkg"

	"github.com/a-novel/service-story-schematics/config"
)

var globalSigner = mustJsonKeysClient()

func mustJsonKeysClient() *jkPkg.ClaimsSigner {
	jsonKeysClient, err := jkPkg.NewAPIClient(context.Background(), config.API.Dependencies.JSONKeys.URL)
	if err != nil {
		panic(err.Error())
	}

	signer := jkPkg.NewClaimsSigner(jsonKeysClient)

	return signer
}

func mustAccessToken(claims authModels.AccessTokenClaims) string {
	accessToken, err := globalSigner.SignClaims(
		context.Background(),
		jkModels.KeyUsageAuth,
		claims,
	)
	if err != nil {
		panic("failed to sign access token: " + err.Error())
	}

	return accessToken
}
