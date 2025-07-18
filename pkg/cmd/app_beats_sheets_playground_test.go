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

		rawideas, err := client.GenerateLoglines(t.Context(), &apimodels.GenerateLoglinesForm{
			Count: 1,
			Theme: "scifi, like Asimov Foundation",
			Lang:  apimodels.LangEn,
		})
		require.NoError(t, err)

		ideas, ok := rawideas.(*apimodels.GenerateLoglinesOKApplicationJSON)
		require.True(t, ok, rawideas)

		rawLogline, err := client.CreateLogline(t.Context(), &apimodels.CreateLoglineForm{
			Slug:    apimodels.Slug(loglineSlug),
			Name:    (*ideas)[0].Name,
			Content: (*ideas)[0].Content,
			Lang:    apimodels.LangEn,
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*apimodels.Logline)
		require.True(t, ok, rawLogline)

		*logline = *newLogline
	}

	t.Log("CreateStoryPlan")
	{
		security.SetToken(userSuperAdminAccessToken)

		rawStoryPlan, err := client.CreateStoryPlan(t.Context(), &planForm)
		require.NoError(t, err)

		newStoryPlan, ok := rawStoryPlan.(*apimodels.StoryPlan)
		require.True(t, ok, rawStoryPlan)

		*storyPlan = *newStoryPlan
	}

	beatsSheet := new(apimodels.BeatsSheet)

	t.Log("GenerateBeatsSheet")
	{
		security.SetToken(userAnonAccessToken)

		rawRes, err := client.GenerateBeatsSheet(t.Context(), &apimodels.GenerateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
			Lang:        apimodels.LangEn,
		})

		require.NoError(t, err)

		_, ok := rawRes.(*apimodels.ForbiddenError)
		require.True(t, ok, rawRes)

		security.SetToken(userLambdaAccessToken)

		rawGeneratedBeatsSheet, err := client.GenerateBeatsSheet(t.Context(), &apimodels.GenerateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
			Lang:        apimodels.LangEn,
		})
		require.NoError(t, err)

		generatedBeatsSheet, ok := rawGeneratedBeatsSheet.(*apimodels.BeatsSheetIdea)
		require.True(t, ok, rawGeneratedBeatsSheet)

		beatsSheet.Content = generatedBeatsSheet.Content
	}

	t.Log("CreateBeatsSheet")
	{
		security.SetToken(userLambdaAccessToken)

		rawBeatsSheet, err := client.CreateBeatsSheet(t.Context(), &apimodels.CreateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
			Content:     beatsSheet.Content,
			Lang:        apimodels.LangEn,
		})
		require.NoError(t, err)

		newBeatsSheet, ok := rawBeatsSheet.(*apimodels.BeatsSheet)
		require.True(t, ok, rawBeatsSheet)

		require.NotEmpty(t, newBeatsSheet.GetID())
		require.Equal(t, logline.ID, newBeatsSheet.GetLoglineID())
		require.Equal(t, storyPlan.ID, newBeatsSheet.GetStoryPlanID())
		require.Equal(t, beatsSheet.Content, newBeatsSheet.GetContent())

		*beatsSheet = *newBeatsSheet
	}

	t.Log("RegenerateBeats")
	{
		security.SetToken(userLambdaAccessToken)

		rawRegeneratedBeatsSheet, err := client.RegenerateBeats(t.Context(), &apimodels.RegenerateBeatsForm{
			BeatsSheetID:   beatsSheet.ID,
			RegenerateKeys: []string{"themeStated"},
		})
		require.NoError(t, err)

		regeneratedBeatsSheet, ok := rawRegeneratedBeatsSheet.(*apimodels.Beats)
		require.True(t, ok, rawRegeneratedBeatsSheet)

		// Save the new beats sheet.
		rawBeatsSheet, err := client.CreateBeatsSheet(t.Context(), &apimodels.CreateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
			Content:     *regeneratedBeatsSheet,
			Lang:        apimodels.LangEn,
		})
		require.NoError(t, err)

		newBeatsSheet, ok := rawBeatsSheet.(*apimodels.BeatsSheet)
		require.True(t, ok, rawBeatsSheet)

		*beatsSheet = *newBeatsSheet
	}

	t.Log("ExpandBeat")
	{
		security.SetToken(userLambdaAccessToken)

		rawExpandedBeat, err := client.ExpandBeat(t.Context(), &apimodels.ExpandBeatForm{
			BeatsSheetID: beatsSheet.ID,
			TargetKey:    "themeStated",
		})
		require.NoError(t, err)

		expandedBeat, ok := rawExpandedBeat.(*apimodels.Beat)
		require.True(t, ok, rawExpandedBeat)

		require.NotEmpty(t, expandedBeat.GetContent())

		beatsSheet.Content[1] = *expandedBeat
		rawBeatsSheet, err := client.CreateBeatsSheet(t.Context(), &apimodels.CreateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
			Content:     beatsSheet.Content,
			Lang:        apimodels.LangEn,
		})
		require.NoError(t, err)

		newBeatsSheet, ok := rawBeatsSheet.(*apimodels.BeatsSheet)
		require.True(t, ok, rawBeatsSheet)

		*beatsSheet = *newBeatsSheet
	}

	t.Log("GetBeatsSheet")
	{
		security.SetToken(userLambdaAccessToken)

		rawBeatsSheet, err := client.GetBeatsSheet(t.Context(), apimodels.GetBeatsSheetParams{
			BeatsSheetID: beatsSheet.ID,
		})

		require.NoError(t, err)

		newBeatsSheet, ok := rawBeatsSheet.(*apimodels.BeatsSheet)
		require.True(t, ok, rawBeatsSheet)

		require.Equal(t, beatsSheet, newBeatsSheet)
	}

	t.Log("ListBeatsSheets")
	{
		security.SetToken(userLambdaAccessToken)

		rawBeatsSheets, err := client.GetBeatsSheets(t.Context(), apimodels.GetBeatsSheetsParams{
			LoglineID: logline.ID,
		})

		require.NoError(t, err)

		beatsSheets, ok := rawBeatsSheets.(*apimodels.GetBeatsSheetsOKApplicationJSON)
		require.True(t, ok, rawBeatsSheets)

		require.Len(t, *beatsSheets, 3)
		require.Equal(t, apimodels.BeatsSheetPreview{
			ID:        beatsSheet.ID,
			Lang:      beatsSheet.Lang,
			CreatedAt: beatsSheet.CreatedAt,
		}, (*beatsSheets)[0])
	}
}
