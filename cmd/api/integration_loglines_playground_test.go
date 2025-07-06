package main

import (
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	authModels "github.com/a-novel/service-authentication/models"

	"github.com/a-novel/service-story-schematics/internal/api/codegen"
)

func TestLoglinesPlayground(t *testing.T) {
	client, securityClient, err := getServerClient()
	require.NoError(t, err)

	loglineSlug := "logline-playground-integration-test"

	loglineIdea := new(codegen.LoglineIdea)

	loglines := make([]*codegen.Logline, 0)

	userLambdaClaims := authModels.AccessTokenClaims{
		UserID: lo.ToPtr(uuid.New()),
		Roles:  []authModels.Role{authModels.RoleUser},
	}
	userLambda2Claims := authModels.AccessTokenClaims{
		UserID: lo.ToPtr(uuid.New()),
		Roles:  []authModels.Role{authModels.RoleUser},
	}
	userAnonClaims := authModels.AccessTokenClaims{
		Roles: []authModels.Role{authModels.RoleAnon},
	}

	userLambdaAccessToken := mustAccessToken(userLambdaClaims)
	userLambda2AccessToken := mustAccessToken(userLambda2Claims)
	userAnonAccessToken := mustAccessToken(userAnonClaims)

	t.Log("LoglineIdeas")
	{
		securityClient.SetToken(userLambdaAccessToken)

		rawideas, err := client.GenerateLoglines(t.Context(), &codegen.GenerateLoglinesForm{
			Count: 2,
			Theme: "scifi, like Asimov Foundation",
			Lang:  codegen.LangEn,
		})
		require.NoError(t, err)

		ideas, ok := rawideas.(*codegen.GenerateLoglinesOKApplicationJSON)
		require.True(t, ok, rawideas)

		require.Len(t, *ideas, 2)

		*loglineIdea = (*ideas)[0]
	}

	t.Log("ExpandLogline")
	{
		securityClient.SetToken(userLambdaAccessToken)

		rawExpandedIdea, err := client.ExpandLogline(t.Context(), loglineIdea)
		require.NoError(t, err)

		expandedIdea, ok := rawExpandedIdea.(*codegen.LoglineIdea)
		require.True(t, ok, rawExpandedIdea)

		*loglineIdea = *expandedIdea
	}

	t.Log("CreateLoglineNotAllowed")
	{
		securityClient.SetToken(userAnonAccessToken)

		rawRes, err := client.CreateLogline(t.Context(), &codegen.CreateLoglineForm{
			Slug:    codegen.Slug(loglineSlug),
			Name:    loglineIdea.Name,
			Content: loglineIdea.Content,
			Lang:    codegen.LangEn,
		})

		require.NoError(t, err)

		_, ok := rawRes.(*codegen.ForbiddenError)
		require.True(t, ok, rawRes)
	}

	t.Log("CreateLogline")
	{
		securityClient.SetToken(userLambdaAccessToken)

		rawLogline, err := client.CreateLogline(t.Context(), &codegen.CreateLoglineForm{
			Slug:    codegen.Slug(loglineSlug),
			Name:    loglineIdea.Name,
			Content: loglineIdea.Content,
			Lang:    codegen.LangEn,
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*codegen.Logline)
		require.True(t, ok, rawLogline)

		require.Equal(t, loglineIdea.Name, newLogline.Name)
		require.Equal(t, loglineIdea.Content, newLogline.Content)
		require.Equal(t, codegen.Slug(loglineSlug), newLogline.Slug)
		require.Equal(t, codegen.UserID(*userLambdaClaims.UserID), newLogline.UserID)

		loglines = append(loglines, newLogline)
	}

	t.Log("GetLoglineByID")
	{
		securityClient.SetToken(userLambdaAccessToken)

		rawLogline, err := client.GetLogline(t.Context(), codegen.GetLoglineParams{
			ID: codegen.OptLoglineID{Value: loglines[0].ID, Set: true},
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*codegen.Logline)
		require.True(t, ok, rawLogline)

		require.Equal(t, loglines[0], newLogline)
	}

	t.Log("GetLoglineBySlug")
	{
		securityClient.SetToken(userLambdaAccessToken)

		rawLogline, err := client.GetLogline(t.Context(), codegen.GetLoglineParams{
			Slug: codegen.OptSlug{Value: loglines[0].Slug, Set: true},
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*codegen.Logline)
		require.True(t, ok, rawLogline)

		require.Equal(t, loglines[0], newLogline)
	}

	t.Log("CreateLogline/SlugResolution")
	{
		securityClient.SetToken(userLambdaAccessToken)

		rawLogline, err := client.CreateLogline(t.Context(), &codegen.CreateLoglineForm{
			Slug:    codegen.Slug(loglineSlug),
			Name:    loglineIdea.Name + " Alt",
			Content: loglineIdea.Content + " Alt",
			Lang:    codegen.LangEn,
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*codegen.Logline)
		require.True(t, ok, rawLogline)

		require.Equal(t, loglineIdea.Name+" Alt", newLogline.Name)
		require.Equal(t, loglineIdea.Content+" Alt", newLogline.Content)
		require.Equal(t, codegen.Slug(loglineSlug+"-1"), newLogline.Slug)
		require.Equal(t, codegen.UserID(*userLambdaClaims.UserID), newLogline.UserID)

		loglines = append(loglines, newLogline)
	}

	t.Log("ListLoglines")
	{
		securityClient.SetToken(userLambdaAccessToken)

		rawLoglines, err := client.GetLoglines(t.Context(), codegen.GetLoglinesParams{})
		require.NoError(t, err)

		userLoglines, ok := rawLoglines.(*codegen.GetLoglinesOKApplicationJSON)
		require.True(t, ok, rawLoglines)

		require.Equal(t, &codegen.GetLoglinesOKApplicationJSON{
			{
				Slug:      loglines[1].Slug,
				Name:      loglines[1].Name,
				Content:   loglines[1].Content,
				Lang:      codegen.LangEn,
				CreatedAt: loglines[1].CreatedAt,
			},
			{
				Slug:      loglines[0].Slug,
				Name:      loglines[0].Name,
				Content:   loglines[0].Content,
				Lang:      codegen.LangEn,
				CreatedAt: loglines[0].CreatedAt,
			},
		}, userLoglines)
	}

	t.Log("NewUserLogline")
	{
		securityClient.SetToken(userLambda2AccessToken)

		rawLogline, err := client.CreateLogline(t.Context(), &codegen.CreateLoglineForm{
			Slug:    codegen.Slug(loglineSlug),
			Name:    loglineIdea.Name,
			Content: loglineIdea.Content,
			Lang:    codegen.LangEn,
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*codegen.Logline)
		require.True(t, ok, rawLogline)

		require.Equal(t, loglineIdea.Name, newLogline.Name)
		require.Equal(t, loglineIdea.Content, newLogline.Content)
		require.Equal(t, codegen.Slug(loglineSlug), newLogline.Slug)
		require.Equal(t, codegen.UserID(*userLambda2Claims.UserID), newLogline.UserID)

		loglines = append(loglines, newLogline)
	}

	t.Log("ListOnlyUserLoglines")
	{
		securityClient.SetToken(userLambda2AccessToken)

		rawLoglines, err := client.GetLoglines(t.Context(), codegen.GetLoglinesParams{})
		require.NoError(t, err)

		userLoglines, ok := rawLoglines.(*codegen.GetLoglinesOKApplicationJSON)
		require.True(t, ok, rawLoglines)

		require.Equal(t, &codegen.GetLoglinesOKApplicationJSON{
			{
				Slug:      loglines[2].Slug,
				Name:      loglines[2].Name,
				Content:   loglines[2].Content,
				CreatedAt: loglines[2].CreatedAt,
				Lang:      codegen.LangEn,
			},
		}, userLoglines)
	}
}
