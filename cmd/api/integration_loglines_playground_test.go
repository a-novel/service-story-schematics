package main

import (
	"crypto/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/authentication/api/apiclient/testapiclient"
	authmodels "github.com/a-novel/authentication/models"

	"github.com/a-novel/story-schematics/api/codegen"
)

func TestLoglinesPlayground(t *testing.T) {
	client, securityClient, err := getServerClient()
	require.NoError(t, err)

	userLambda := rand.Text()
	userLambdaID := uuid.New()

	testapiclient.AddPool(userLambda, &authmodels.AccessTokenClaims{
		UserID: &userLambdaID,
		Roles:  []authmodels.Role{authmodels.RoleUser},
	})

	userAnon := rand.Text()
	testapiclient.AddPool(userAnon, &authmodels.AccessTokenClaims{})

	loglineSlug := "logline-playground-integration-test"

	loglineIdea := new(codegen.LoglineIdea)

	loglines := make([]*codegen.Logline, 0)

	t.Log("LoglineIdeas")
	{
		securityClient.SetToken(userLambda)

		rawideas, err := client.GenerateLoglines(t.Context(), &codegen.GenerateLoglinesForm{
			Count: 2,
			Theme: "scifi, like Asimov Foundation",
		})
		require.NoError(t, err)

		ideas, ok := rawideas.(*codegen.GenerateLoglinesOKApplicationJSON)
		require.True(t, ok)

		require.Len(t, *ideas, 2)

		*loglineIdea = (*ideas)[0]
	}

	t.Log("Expandlogline")
	{
		securityClient.SetToken(userLambda)

		rawExpandedIdea, err := client.ExpandLogline(t.Context(), loglineIdea)
		require.NoError(t, err)

		expandedIdea, ok := rawExpandedIdea.(*codegen.LoglineIdea)
		require.True(t, ok)

		*loglineIdea = *expandedIdea
	}

	t.Log("CreateLoglineNotAllowed")
	{
		securityClient.SetToken(userAnon)

		_, err = client.CreateLogline(t.Context(), &codegen.CreateLoglineForm{
			Slug:    codegen.Slug(loglineSlug),
			Name:    loglineIdea.Name,
			Content: loglineIdea.Content,
		})
		require.Error(t, err)
	}

	t.Log("CreateLogline")
	{
		securityClient.SetToken(userLambda)

		rawLogline, err := client.CreateLogline(t.Context(), &codegen.CreateLoglineForm{
			Slug:    codegen.Slug(loglineSlug),
			Name:    loglineIdea.Name,
			Content: loglineIdea.Content,
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*codegen.Logline)
		require.True(t, ok)

		require.Equal(t, loglineIdea.Name, newLogline.Name)
		require.Equal(t, loglineIdea.Content, newLogline.Content)
		require.Equal(t, codegen.Slug(loglineSlug), newLogline.Slug)
		require.Equal(t, codegen.UserID(userLambdaID), newLogline.UserID)

		loglines = append(loglines, newLogline)
	}

	t.Log("GetLoglineByID")
	{
		securityClient.SetToken(userLambda)

		rawLogline, err := client.GetLogline(t.Context(), codegen.GetLoglineParams{
			ID: codegen.OptLoglineID{Value: loglines[0].ID, Set: true},
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*codegen.Logline)
		require.True(t, ok)

		require.Equal(t, loglines[0], newLogline)
	}

	t.Log("GetLoglineBySlug")
	{
		securityClient.SetToken(userLambda)

		rawLogline, err := client.GetLogline(t.Context(), codegen.GetLoglineParams{
			Slug: codegen.OptSlug{Value: loglines[0].Slug, Set: true},
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*codegen.Logline)
		require.True(t, ok)

		require.Equal(t, loglines[0], newLogline)
	}

	t.Log("CreateLogline/SlugResolution")
	{
		securityClient.SetToken(userLambda)

		rawLogline, err := client.CreateLogline(t.Context(), &codegen.CreateLoglineForm{
			Slug:    codegen.Slug(loglineSlug),
			Name:    loglineIdea.Name + " Alt",
			Content: loglineIdea.Content + " Alt",
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*codegen.Logline)
		require.True(t, ok)

		require.Equal(t, loglineIdea.Name+" Alt", newLogline.Name)
		require.Equal(t, loglineIdea.Content+" Alt", newLogline.Content)
		require.Equal(t, codegen.Slug(loglineSlug+"-1"), newLogline.Slug)
		require.Equal(t, codegen.UserID(userLambdaID), newLogline.UserID)

		loglines = append(loglines, newLogline)
	}

	t.Log("ListLoglines")
	{
		securityClient.SetToken(userLambda)

		rawLoglines, err := client.GetLoglines(t.Context(), codegen.GetLoglinesParams{})
		require.NoError(t, err)

		userLoglines, ok := rawLoglines.(*codegen.GetLoglinesOKApplicationJSON)
		require.True(t, ok)

		require.Equal(t, &codegen.GetLoglinesOKApplicationJSON{
			{
				Slug:      loglines[1].Slug,
				Name:      loglines[1].Name,
				Content:   loglines[1].Content,
				CreatedAt: loglines[1].CreatedAt,
			},
			{
				Slug:      loglines[0].Slug,
				Name:      loglines[0].Name,
				Content:   loglines[0].Content,
				CreatedAt: loglines[0].CreatedAt,
			},
		}, userLoglines)
	}

	userLambda2 := rand.Text()
	userLambda2ID := uuid.New()

	testapiclient.AddPool(userLambda2, &authmodels.AccessTokenClaims{
		UserID: &userLambda2ID,
		Roles:  []authmodels.Role{authmodels.RoleUser},
	})

	t.Log("NewUserLogline")
	{
		securityClient.SetToken(userLambda2)

		rawLogline, err := client.CreateLogline(t.Context(), &codegen.CreateLoglineForm{
			Slug:    codegen.Slug(loglineSlug),
			Name:    loglineIdea.Name,
			Content: loglineIdea.Content,
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*codegen.Logline)
		require.True(t, ok)

		require.Equal(t, loglineIdea.Name, newLogline.Name)
		require.Equal(t, loglineIdea.Content, newLogline.Content)
		require.Equal(t, codegen.Slug(loglineSlug), newLogline.Slug)
		require.Equal(t, codegen.UserID(userLambda2ID), newLogline.UserID)

		loglines = append(loglines, newLogline)
	}

	t.Log("ListOnlyUserLoglines")
	{
		securityClient.SetToken(userLambda2)

		rawLoglines, err := client.GetLoglines(t.Context(), codegen.GetLoglinesParams{})
		require.NoError(t, err)

		userLoglines, ok := rawLoglines.(*codegen.GetLoglinesOKApplicationJSON)
		require.True(t, ok)

		require.Equal(t, &codegen.GetLoglinesOKApplicationJSON{
			{
				Slug:      loglines[2].Slug,
				Name:      loglines[2].Name,
				Content:   loglines[2].Content,
				CreatedAt: loglines[2].CreatedAt,
			},
		}, userLoglines)
	}
}
