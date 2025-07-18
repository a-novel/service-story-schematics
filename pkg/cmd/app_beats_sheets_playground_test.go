package cmdpkg_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/ogen"
	authmodels "github.com/a-novel/service-authentication/models"

	apimodels "github.com/a-novel/service-story-schematics/models/api"
	"github.com/a-novel/service-story-schematics/pkg"
)

func testAppBeatsSheetsPlayground(ctx context.Context, t *testing.T, appConfig TestConfig) {
	t.Helper()

	security := pkg.NewBearerSource()

	client, err := pkg.NewAPIClient(ctx, fmt.Sprintf("http://localhost:%v/v1", appConfig.API.Port), security)
	require.NoError(t, err)

	loglineSlug := "beats-sheets-playground-integration-test"

	storyPlanSlug := "beats-sheets-playground-integration-test-save-the-cat-partial"
	planForm := *saveTheCatPartialPlanForm
	planForm.Slug = apimodels.Slug(storyPlanSlug)

	logline := new(apimodels.Logline)
	storyPlan := new(apimodels.StoryPlan)

	userLambdaClaims := authmodels.AccessTokenClaims{
		UserID: lo.ToPtr(uuid.New()),
		Roles:  []authmodels.Role{authmodels.RoleUser},
	}
	userSuperAdminClaims := authmodels.AccessTokenClaims{
		UserID: lo.ToPtr(uuid.New()),
		Roles:  []authmodels.Role{authmodels.RoleSuperAdmin},
	}
	userAnonClaims := authmodels.AccessTokenClaims{
		Roles: []authmodels.Role{authmodels.RoleAnon},
	}

	userLambdaAccessToken := getAccessToken(t, appConfig, userLambdaClaims)
	userSuperAdminAccessToken := getAccessToken(t, appConfig, userSuperAdminClaims)
	userAnonAccessToken := getAccessToken(t, appConfig, userAnonClaims)

	t.Log("CreateLogline")
	{
		security.SetToken(userLambdaAccessToken)

		ideas, err := ogen.MustGetResponse[
			apimodels.GenerateLoglinesRes, *apimodels.GenerateLoglinesOKApplicationJSON,
		](
			client.GenerateLoglines(t.Context(), &apimodels.GenerateLoglinesForm{
				Count: 1,
				Theme: "scifi, like Asimov Foundation",
				Lang:  apimodels.LangEn,
			}),
		)
		require.NoError(t, err)

		newLogline, err := ogen.MustGetResponse[apimodels.CreateLoglineRes, *apimodels.Logline](
			client.CreateLogline(t.Context(), &apimodels.CreateLoglineForm{
				Slug:    apimodels.Slug(loglineSlug),
				Name:    (*ideas)[0].Name,
				Content: (*ideas)[0].Content,
				Lang:    apimodels.LangEn,
			}),
		)
		require.NoError(t, err)

		*logline = *newLogline
	}

	t.Log("CreateStoryPlan")
	{
		security.SetToken(userSuperAdminAccessToken)

		newStoryPlan, err := ogen.MustGetResponse[apimodels.CreateStoryPlanRes, *apimodels.StoryPlan](
			client.CreateStoryPlan(t.Context(), &planForm),
		)
		require.NoError(t, err)

		*storyPlan = *newStoryPlan
	}

	beatsSheet := new(apimodels.BeatsSheet)

	t.Log("GenerateBeatsSheet")
	{
		security.SetToken(userAnonAccessToken)

		_, err = ogen.MustGetResponse[apimodels.GenerateBeatsSheetRes, *apimodels.ForbiddenError](
			client.GenerateBeatsSheet(t.Context(), &apimodels.GenerateBeatsSheetForm{
				LoglineID:   logline.ID,
				StoryPlanID: storyPlan.ID,
				Lang:        apimodels.LangEn,
			}),
		)
		require.NoError(t, err)

		security.SetToken(userLambdaAccessToken)

		generatedBeatsSheet, err := ogen.MustGetResponse[apimodels.GenerateBeatsSheetRes, *apimodels.BeatsSheetIdea](
			client.GenerateBeatsSheet(t.Context(), &apimodels.GenerateBeatsSheetForm{
				LoglineID:   logline.ID,
				StoryPlanID: storyPlan.ID,
				Lang:        apimodels.LangEn,
			}),
		)
		require.NoError(t, err)

		beatsSheet.Content = generatedBeatsSheet.Content
	}

	t.Log("CreateBeatsSheet")
	{
		security.SetToken(userLambdaAccessToken)

		newBeatsSheet, err := ogen.MustGetResponse[apimodels.CreateBeatsSheetRes, *apimodels.BeatsSheet](
			client.CreateBeatsSheet(t.Context(), &apimodels.CreateBeatsSheetForm{
				LoglineID:   logline.ID,
				StoryPlanID: storyPlan.ID,
				Content:     beatsSheet.Content,
				Lang:        apimodels.LangEn,
			}),
		)
		require.NoError(t, err)

		require.NotEmpty(t, newBeatsSheet.GetID())
		require.Equal(t, logline.ID, newBeatsSheet.GetLoglineID())
		require.Equal(t, storyPlan.ID, newBeatsSheet.GetStoryPlanID())
		require.Equal(t, beatsSheet.Content, newBeatsSheet.GetContent())

		*beatsSheet = *newBeatsSheet
	}

	t.Log("RegenerateBeats")
	{
		security.SetToken(userLambdaAccessToken)

		regeneratedBeatsSheet, err := ogen.MustGetResponse[apimodels.RegenerateBeatsRes, *apimodels.Beats](
			client.RegenerateBeats(t.Context(), &apimodels.RegenerateBeatsForm{
				BeatsSheetID:   beatsSheet.ID,
				RegenerateKeys: []string{"themeStated"},
			}),
		)
		require.NoError(t, err)

		newBeatsSheet, err := ogen.MustGetResponse[apimodels.CreateBeatsSheetRes, *apimodels.BeatsSheet](
			client.CreateBeatsSheet(t.Context(), &apimodels.CreateBeatsSheetForm{
				LoglineID:   logline.ID,
				StoryPlanID: storyPlan.ID,
				Content:     *regeneratedBeatsSheet,
				Lang:        apimodels.LangEn,
			}),
		)
		require.NoError(t, err)

		*beatsSheet = *newBeatsSheet
	}

	t.Log("ExpandBeat")
	{
		security.SetToken(userLambdaAccessToken)

		expandedBeat, err := ogen.MustGetResponse[apimodels.ExpandBeatRes, *apimodels.Beat](
			client.ExpandBeat(t.Context(), &apimodels.ExpandBeatForm{
				BeatsSheetID: beatsSheet.ID,
				TargetKey:    "themeStated",
			}),
		)
		require.NoError(t, err)

		require.NotEmpty(t, expandedBeat.GetContent())

		beatsSheet.Content[1] = *expandedBeat

		newBeatsSheet, err := ogen.MustGetResponse[apimodels.CreateBeatsSheetRes, *apimodels.BeatsSheet](
			client.CreateBeatsSheet(t.Context(), &apimodels.CreateBeatsSheetForm{
				LoglineID:   logline.ID,
				StoryPlanID: storyPlan.ID,
				Content:     beatsSheet.Content,
				Lang:        apimodels.LangEn,
			}),
		)
		require.NoError(t, err)

		*beatsSheet = *newBeatsSheet
	}

	t.Log("GetBeatsSheet")
	{
		security.SetToken(userLambdaAccessToken)

		newBeatsSheet, err := ogen.MustGetResponse[apimodels.GetBeatsSheetRes, *apimodels.BeatsSheet](
			client.GetBeatsSheet(t.Context(), apimodels.GetBeatsSheetParams{BeatsSheetID: beatsSheet.ID}),
		)
		require.NoError(t, err)

		require.Equal(t, beatsSheet, newBeatsSheet)
	}

	t.Log("ListBeatsSheets")
	{
		security.SetToken(userLambdaAccessToken)

		beatsSheets, err := ogen.MustGetResponse[
			apimodels.GetBeatsSheetsRes, *apimodels.GetBeatsSheetsOKApplicationJSON,
		](
			client.GetBeatsSheets(t.Context(), apimodels.GetBeatsSheetsParams{LoglineID: logline.ID}),
		)
		require.NoError(t, err)

		require.Len(t, *beatsSheets, 3)
		require.Equal(t, apimodels.BeatsSheetPreview{
			ID:        beatsSheet.ID,
			Lang:      beatsSheet.Lang,
			CreatedAt: beatsSheet.CreatedAt,
		}, (*beatsSheets)[0])
	}
}
