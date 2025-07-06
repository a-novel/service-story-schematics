package main

import (
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	authModels "github.com/a-novel/service-authentication/models"

	"github.com/a-novel/service-story-schematics/internal/api/codegen"
)

func TestBeatsSheetsPlayground(t *testing.T) {
	client, securityClient, err := getServerClient()
	require.NoError(t, err)

	loglineSlug := "beats-sheets-playground-integration-test"

	storyPlanSlug := "beats-sheets-playground-integration-test-save-the-cat-partial"
	planForm := *saveTheCatPartialPlanForm
	planForm.Slug = codegen.Slug(storyPlanSlug)

	logline := new(codegen.Logline)
	storyPlan := new(codegen.StoryPlan)

	userLambdaClaims := authModels.AccessTokenClaims{
		UserID: lo.ToPtr(uuid.New()),
		Roles:  []authModels.Role{authModels.RoleUser},
	}
	userSuperAdminClaims := authModels.AccessTokenClaims{
		UserID: lo.ToPtr(uuid.New()),
		Roles:  []authModels.Role{authModels.RoleSuperAdmin},
	}
	userAnonClaims := authModels.AccessTokenClaims{
		Roles: []authModels.Role{authModels.RoleAnon},
	}

	userLambdaAccessToken := mustAccessToken(userLambdaClaims)
	userSuperAdminAccessToken := mustAccessToken(userSuperAdminClaims)
	userAnonAccessToken := mustAccessToken(userAnonClaims)

	t.Log("CreateLogline")
	{
		securityClient.SetToken(userLambdaAccessToken)

		rawideas, err := client.GenerateLoglines(t.Context(), &codegen.GenerateLoglinesForm{
			Count: 1,
			Theme: "scifi, like Asimov Foundation",
			Lang:  codegen.LangEn,
		})
		require.NoError(t, err)

		ideas, ok := rawideas.(*codegen.GenerateLoglinesOKApplicationJSON)
		require.True(t, ok, rawideas)

		rawLogline, err := client.CreateLogline(t.Context(), &codegen.CreateLoglineForm{
			Slug:    codegen.Slug(loglineSlug),
			Name:    (*ideas)[0].Name,
			Content: (*ideas)[0].Content,
			Lang:    codegen.LangEn,
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*codegen.Logline)
		require.True(t, ok, rawLogline)

		*logline = *newLogline
	}

	t.Log("CreateStoryPlan")
	{
		securityClient.SetToken(userSuperAdminAccessToken)

		rawStoryPlan, err := client.CreateStoryPlan(t.Context(), &planForm)
		require.NoError(t, err)

		newStoryPlan, ok := rawStoryPlan.(*codegen.StoryPlan)
		require.True(t, ok, rawStoryPlan)

		*storyPlan = *newStoryPlan
	}

	beatsSheet := new(codegen.BeatsSheet)

	t.Log("GenerateBeatsSheet")
	{
		securityClient.SetToken(userAnonAccessToken)

		rawRes, err := client.GenerateBeatsSheet(t.Context(), &codegen.GenerateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
			Lang:        codegen.LangEn,
		})

		require.NoError(t, err)

		_, ok := rawRes.(*codegen.ForbiddenError)
		require.True(t, ok, rawRes)

		securityClient.SetToken(userLambdaAccessToken)

		rawGeneratedBeatsSheet, err := client.GenerateBeatsSheet(t.Context(), &codegen.GenerateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
			Lang:        codegen.LangEn,
		})
		require.NoError(t, err)

		generatedBeatsSheet, ok := rawGeneratedBeatsSheet.(*codegen.BeatsSheetIdea)
		require.True(t, ok, rawGeneratedBeatsSheet)

		beatsSheet.Content = generatedBeatsSheet.Content
	}

	t.Log("CreateBeatsSheet")
	{
		securityClient.SetToken(userLambdaAccessToken)

		rawBeatsSheet, err := client.CreateBeatsSheet(t.Context(), &codegen.CreateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
			Content:     beatsSheet.Content,
			Lang:        codegen.LangEn,
		})
		require.NoError(t, err)

		newBeatsSheet, ok := rawBeatsSheet.(*codegen.BeatsSheet)
		require.True(t, ok, rawBeatsSheet)

		require.NotEmpty(t, newBeatsSheet.GetID())
		require.Equal(t, logline.ID, newBeatsSheet.GetLoglineID())
		require.Equal(t, storyPlan.ID, newBeatsSheet.GetStoryPlanID())
		require.Equal(t, beatsSheet.Content, newBeatsSheet.GetContent())

		*beatsSheet = *newBeatsSheet
	}

	t.Log("RegenerateBeats")
	{
		securityClient.SetToken(userLambdaAccessToken)

		rawRegeneratedBeatsSheet, err := client.RegenerateBeats(t.Context(), &codegen.RegenerateBeatsForm{
			BeatsSheetID:   beatsSheet.ID,
			RegenerateKeys: []string{"themeStated"},
		})
		require.NoError(t, err)

		regeneratedBeatsSheet, ok := rawRegeneratedBeatsSheet.(*codegen.Beats)
		require.True(t, ok, rawRegeneratedBeatsSheet)

		// Save the new beats sheet.
		rawBeatsSheet, err := client.CreateBeatsSheet(t.Context(), &codegen.CreateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
			Content:     *regeneratedBeatsSheet,
			Lang:        codegen.LangEn,
		})
		require.NoError(t, err)

		newBeatsSheet, ok := rawBeatsSheet.(*codegen.BeatsSheet)
		require.True(t, ok, rawBeatsSheet)

		*beatsSheet = *newBeatsSheet
	}

	t.Log("ExpandBeat")
	{
		securityClient.SetToken(userLambdaAccessToken)

		rawExpandedBeat, err := client.ExpandBeat(t.Context(), &codegen.ExpandBeatForm{
			BeatsSheetID: beatsSheet.ID,
			TargetKey:    "themeStated",
		})
		require.NoError(t, err)

		expandedBeat, ok := rawExpandedBeat.(*codegen.Beat)
		require.True(t, ok, rawExpandedBeat)

		require.NotEmpty(t, expandedBeat.GetContent())

		beatsSheet.Content[1] = *expandedBeat
		rawBeatsSheet, err := client.CreateBeatsSheet(t.Context(), &codegen.CreateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
			Content:     beatsSheet.Content,
			Lang:        codegen.LangEn,
		})
		require.NoError(t, err)

		newBeatsSheet, ok := rawBeatsSheet.(*codegen.BeatsSheet)
		require.True(t, ok, rawBeatsSheet)

		*beatsSheet = *newBeatsSheet
	}

	t.Log("GetBeatsSheet")
	{
		securityClient.SetToken(userLambdaAccessToken)

		rawBeatsSheet, err := client.GetBeatsSheet(t.Context(), codegen.GetBeatsSheetParams{
			BeatsSheetID: beatsSheet.ID,
		})

		require.NoError(t, err)

		newBeatsSheet, ok := rawBeatsSheet.(*codegen.BeatsSheet)
		require.True(t, ok, rawBeatsSheet)

		require.Equal(t, beatsSheet, newBeatsSheet)
	}

	t.Log("ListBeatsSheets")
	{
		securityClient.SetToken(userLambdaAccessToken)

		rawBeatsSheets, err := client.GetBeatsSheets(t.Context(), codegen.GetBeatsSheetsParams{
			LoglineID: logline.ID,
		})

		require.NoError(t, err)

		beatsSheets, ok := rawBeatsSheets.(*codegen.GetBeatsSheetsOKApplicationJSON)
		require.True(t, ok, rawBeatsSheets)

		require.Len(t, *beatsSheets, 3)
		require.Equal(t, codegen.BeatsSheetPreview{
			ID:        beatsSheet.ID,
			Lang:      beatsSheet.Lang,
			CreatedAt: beatsSheet.CreatedAt,
		}, (*beatsSheets)[0])
	}
}
