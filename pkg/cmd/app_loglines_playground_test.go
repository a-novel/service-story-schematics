package cmdpkg_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	authmodels "github.com/a-novel/service-authentication/models"

	apimodels "github.com/a-novel/service-story-schematics/models/api"
	"github.com/a-novel/service-story-schematics/pkg"
)

func testAppLoglinesPlayground(ctx context.Context, t *testing.T, appConfig TestConfig) {
	t.Helper()

	security := pkg.NewBearerSource()

	client, err := pkg.NewAPIClient(ctx, fmt.Sprintf("http://localhost:%v/v1", appConfig.API.Port), security)
	require.NoError(t, err)

	loglineSlug := "logline-playground-integration-test"

	loglineIdea := new(apimodels.LoglineIdea)

	loglines := make([]*apimodels.Logline, 0)

	userLambdaClaims := authmodels.AccessTokenClaims{
		UserID: lo.ToPtr(uuid.New()),
		Roles:  []authmodels.Role{authmodels.RoleUser},
	}
	userLambda2Claims := authmodels.AccessTokenClaims{
		UserID: lo.ToPtr(uuid.New()),
		Roles:  []authmodels.Role{authmodels.RoleUser},
	}
	userAnonClaims := authmodels.AccessTokenClaims{
		Roles: []authmodels.Role{authmodels.RoleAnon},
	}

	userLambdaAccessToken := getAccessToken(t, appConfig, userLambdaClaims)
	userLambda2AccessToken := getAccessToken(t, appConfig, userLambda2Claims)
	userAnonAccessToken := getAccessToken(t, appConfig, userAnonClaims)

	t.Log("LoglineIdeas")
	{
		security.SetToken(userLambdaAccessToken)

		rawideas, err := client.GenerateLoglines(t.Context(), &apimodels.GenerateLoglinesForm{
			Count: 2,
			Theme: "scifi, like Asimov Foundation",
			Lang:  apimodels.LangEn,
		})
		require.NoError(t, err)

		ideas, ok := rawideas.(*apimodels.GenerateLoglinesOKApplicationJSON)
		require.True(t, ok, rawideas)

		require.Len(t, *ideas, 2)

		*loglineIdea = (*ideas)[0]
	}

	t.Log("ExpandLogline")
	{
		security.SetToken(userLambdaAccessToken)

		rawExpandedIdea, err := client.ExpandLogline(t.Context(), loglineIdea)
		require.NoError(t, err)

		expandedIdea, ok := rawExpandedIdea.(*apimodels.LoglineIdea)
		require.True(t, ok, rawExpandedIdea)

		*loglineIdea = *expandedIdea
	}

	t.Log("CreateLoglineNotAllowed")
	{
		security.SetToken(userAnonAccessToken)

		rawRes, err := client.CreateLogline(t.Context(), &apimodels.CreateLoglineForm{
			Slug:    apimodels.Slug(loglineSlug),
			Name:    loglineIdea.Name,
			Content: loglineIdea.Content,
			Lang:    apimodels.LangEn,
		})

		require.NoError(t, err)

		_, ok := rawRes.(*apimodels.ForbiddenError)
		require.True(t, ok, rawRes)
	}

	t.Log("CreateLogline")
	{
		security.SetToken(userLambdaAccessToken)

		rawLogline, err := client.CreateLogline(t.Context(), &apimodels.CreateLoglineForm{
			Slug:    apimodels.Slug(loglineSlug),
			Name:    loglineIdea.Name,
			Content: loglineIdea.Content,
			Lang:    apimodels.LangEn,
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*apimodels.Logline)
		require.True(t, ok, rawLogline)

		require.Equal(t, loglineIdea.Name, newLogline.Name)
		require.Equal(t, loglineIdea.Content, newLogline.Content)
		require.Equal(t, apimodels.Slug(loglineSlug), newLogline.Slug)
		require.Equal(t, apimodels.UserID(*userLambdaClaims.UserID), newLogline.UserID)

		loglines = append(loglines, newLogline)
	}

	t.Log("GetLoglineByID")
	{
		security.SetToken(userLambdaAccessToken)

		rawLogline, err := client.GetLogline(t.Context(), apimodels.GetLoglineParams{
			ID: apimodels.OptLoglineID{Value: loglines[0].ID, Set: true},
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*apimodels.Logline)
		require.True(t, ok, rawLogline)

		require.Equal(t, loglines[0], newLogline)
	}

	t.Log("GetLoglineBySlug")
	{
		security.SetToken(userLambdaAccessToken)

		rawLogline, err := client.GetLogline(t.Context(), apimodels.GetLoglineParams{
			Slug: apimodels.OptSlug{Value: loglines[0].Slug, Set: true},
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*apimodels.Logline)
		require.True(t, ok, rawLogline)

		require.Equal(t, loglines[0], newLogline)
	}

	t.Log("CreateLogline/SlugResolution")
	{
		security.SetToken(userLambdaAccessToken)

		rawLogline, err := client.CreateLogline(t.Context(), &apimodels.CreateLoglineForm{
			Slug:    apimodels.Slug(loglineSlug),
			Name:    loglineIdea.Name + " Alt",
			Content: loglineIdea.Content + " Alt",
			Lang:    apimodels.LangEn,
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*apimodels.Logline)
		require.True(t, ok, rawLogline)

		require.Equal(t, loglineIdea.Name+" Alt", newLogline.Name)
		require.Equal(t, loglineIdea.Content+" Alt", newLogline.Content)
		require.Equal(t, apimodels.Slug(loglineSlug+"-1"), newLogline.Slug)
		require.Equal(t, apimodels.UserID(*userLambdaClaims.UserID), newLogline.UserID)

		loglines = append(loglines, newLogline)
	}

	t.Log("ListLoglines")
	{
		security.SetToken(userLambdaAccessToken)

		rawLoglines, err := client.GetLoglines(t.Context(), apimodels.GetLoglinesParams{})
		require.NoError(t, err)

		userLoglines, ok := rawLoglines.(*apimodels.GetLoglinesOKApplicationJSON)
		require.True(t, ok, rawLoglines)

		require.Equal(t, &apimodels.GetLoglinesOKApplicationJSON{
			{
				Slug:      loglines[1].Slug,
				Name:      loglines[1].Name,
				Content:   loglines[1].Content,
				Lang:      apimodels.LangEn,
				CreatedAt: loglines[1].CreatedAt,
			},
			{
				Slug:      loglines[0].Slug,
				Name:      loglines[0].Name,
				Content:   loglines[0].Content,
				Lang:      apimodels.LangEn,
				CreatedAt: loglines[0].CreatedAt,
			},
		}, userLoglines)
	}

	t.Log("NewUserLogline")
	{
		security.SetToken(userLambda2AccessToken)

		rawLogline, err := client.CreateLogline(t.Context(), &apimodels.CreateLoglineForm{
			Slug:    apimodels.Slug(loglineSlug),
			Name:    loglineIdea.Name,
			Content: loglineIdea.Content,
			Lang:    apimodels.LangEn,
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*apimodels.Logline)
		require.True(t, ok, rawLogline)

		require.Equal(t, loglineIdea.Name, newLogline.Name)
		require.Equal(t, loglineIdea.Content, newLogline.Content)
		require.Equal(t, apimodels.Slug(loglineSlug), newLogline.Slug)
		require.Equal(t, apimodels.UserID(*userLambda2Claims.UserID), newLogline.UserID)

		loglines = append(loglines, newLogline)
	}

	t.Log("ListOnlyUserLoglines")
	{
		security.SetToken(userLambda2AccessToken)

		rawLoglines, err := client.GetLoglines(t.Context(), apimodels.GetLoglinesParams{})
		require.NoError(t, err)

		userLoglines, ok := rawLoglines.(*apimodels.GetLoglinesOKApplicationJSON)
		require.True(t, ok, rawLoglines)

		require.Equal(t, &apimodels.GetLoglinesOKApplicationJSON{
			{
				Slug:      loglines[2].Slug,
				Name:      loglines[2].Name,
				Content:   loglines[2].Content,
				CreatedAt: loglines[2].CreatedAt,
				Lang:      apimodels.LangEn,
			},
		}, userLoglines)
	}
}
