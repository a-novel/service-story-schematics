package main

import (
	"crypto/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-authentication/api/apiclient/testapiclient"
	authmodels "github.com/a-novel/service-authentication/models"

	"github.com/a-novel/service-story-schematics/api/codegen"
)

func TestBeastSheetsPlayground(t *testing.T) {
	client, securityClient, err := getServerClient()
	require.NoError(t, err)

	userLambda := rand.Text()
	userLambdaID := uuid.New()
	testapiclient.AddPool(userLambda, &authmodels.AccessTokenClaims{
		UserID: &userLambdaID,
		Roles:  []authmodels.Role{authmodels.RoleUser},
	})

	userSuperAdmin := rand.Text()
	testapiclient.AddPool(userSuperAdmin, &authmodels.AccessTokenClaims{
		UserID: lo.ToPtr(uuid.New()),
		Roles:  []authmodels.Role{authmodels.RoleSuperAdmin},
	})

	userAnon := rand.Text()
	testapiclient.AddPool(userAnon, &authmodels.AccessTokenClaims{})

	loglineSlug := "beats-sheets-playground-integration-test"

	storyPlanSlug := "beats-sheets-playground-integration-test-save-the-cat-partial"
	planForm := *saveTheCatPartialPlanForm
	planForm.Slug = codegen.Slug(storyPlanSlug)

	logline := new(codegen.Logline)
	storyPlan := new(codegen.StoryPlan)

	t.Log("CreateLogline")
	{
		securityClient.SetToken(userLambda)

		rawideas, err := client.GenerateLoglines(t.Context(), &codegen.GenerateLoglinesForm{
			Count: 1,
			Theme: "scifi, like Asimov Foundation",
		})
		require.NoError(t, err)

		ideas, ok := rawideas.(*codegen.GenerateLoglinesOKApplicationJSON)
		require.True(t, ok)

		rawLogline, err := client.CreateLogline(t.Context(), &codegen.CreateLoglineForm{
			Slug:    codegen.Slug(loglineSlug),
			Name:    (*ideas)[0].Name,
			Content: (*ideas)[0].Content,
		})
		require.NoError(t, err)

		newLogline, ok := rawLogline.(*codegen.Logline)
		require.True(t, ok)

		*logline = *newLogline
	}

	t.Log("CreateStoryPlan")
	{
		securityClient.SetToken(userSuperAdmin)

		rawStoryPlan, err := client.CreateStoryPlan(t.Context(), &planForm)
		require.NoError(t, err)

		newStoryPlan, ok := rawStoryPlan.(*codegen.StoryPlan)
		require.True(t, ok)

		*storyPlan = *newStoryPlan
	}

	beatsSheet := new(codegen.BeatsSheet)

	t.Log("GenerateBeatsSheet")
	{
		securityClient.SetToken(userAnon)

		rawRes, err := client.GenerateBeatsSheet(t.Context(), &codegen.GenerateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
		})

		require.NoError(t, err)

		_, ok := rawRes.(*codegen.UnauthorizedError)
		require.True(t, ok)

		securityClient.SetToken(userLambda)

		rawGeneratedBeatsSheet, err := client.GenerateBeatsSheet(t.Context(), &codegen.GenerateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
		})
		require.NoError(t, err)

		generatedBeatsSheet, ok := rawGeneratedBeatsSheet.(*codegen.BeatsSheet)
		require.True(t, ok)

		*beatsSheet = *generatedBeatsSheet
	}

	t.Log("CreateBeatsSheet")
	{
		securityClient.SetToken(userLambda)

		rawBeatsSheet, err := client.CreateBeatsSheet(t.Context(), &codegen.CreateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
			Content:     beatsSheet.Content,
		})
		require.NoError(t, err)

		newBeatsSheet, ok := rawBeatsSheet.(*codegen.BeatsSheet)
		require.True(t, ok)

		require.NotEmpty(t, newBeatsSheet.GetID())
		require.Equal(t, logline.ID, newBeatsSheet.GetLoglineID())
		require.Equal(t, storyPlan.ID, newBeatsSheet.GetStoryPlanID())
		require.Equal(t, beatsSheet.Content, newBeatsSheet.GetContent())

		*beatsSheet = *newBeatsSheet
	}

	t.Log("RegenerateBeats")
	{
		securityClient.SetToken(userLambda)

		rawRegeneratedBeatsSheet, err := client.RegenerateBeats(t.Context(), &codegen.RegenerateBeatsForm{
			BeatsSheetID:   beatsSheet.ID,
			RegenerateKeys: []string{"themeStated"},
		})
		require.NoError(t, err)

		regeneratedBeatsSheet, ok := rawRegeneratedBeatsSheet.(*codegen.BeatsSheet)
		require.True(t, ok)

		// Save the new beats sheet.
		rawBeatsSheet, err := client.CreateBeatsSheet(t.Context(), &codegen.CreateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
			Content:     regeneratedBeatsSheet.Content,
		})
		require.NoError(t, err)

		newBeatsSheet, ok := rawBeatsSheet.(*codegen.BeatsSheet)
		require.True(t, ok)

		*beatsSheet = *newBeatsSheet
	}

	t.Log("ExpandBeat")
	{
		securityClient.SetToken(userLambda)

		rawExpandedBeat, err := client.ExpandBeat(t.Context(), &codegen.ExpandBeatForm{
			BeatsSheetID: beatsSheet.ID,
			TargetKey:    "themeStated",
		})
		require.NoError(t, err)

		expandedBeat, ok := rawExpandedBeat.(*codegen.Beat)
		require.True(t, ok)

		require.NotEmpty(t, expandedBeat.GetContent())

		beatsSheet.Content[1] = *expandedBeat
		rawBeatsSheet, err := client.CreateBeatsSheet(t.Context(), &codegen.CreateBeatsSheetForm{
			LoglineID:   logline.ID,
			StoryPlanID: storyPlan.ID,
			Content:     beatsSheet.Content,
		})
		require.NoError(t, err)

		newBeatsSheet, ok := rawBeatsSheet.(*codegen.BeatsSheet)
		require.True(t, ok)

		*beatsSheet = *newBeatsSheet
	}

	t.Log("GetBeatsSheet")
	{
		securityClient.SetToken(userLambda)

		rawBeatsSheet, err := client.GetBeatsSheet(t.Context(), codegen.GetBeatsSheetParams{
			BeatsSheetID: beatsSheet.ID,
		})

		require.NoError(t, err)

		newBeatsSheet, ok := rawBeatsSheet.(*codegen.BeatsSheet)
		require.True(t, ok)

		require.Equal(t, beatsSheet, newBeatsSheet)
	}

	t.Log("ListBeatsSheets")
	{
		securityClient.SetToken(userLambda)

		rawBeatsSheets, err := client.GetBeatsSheets(t.Context(), codegen.GetBeatsSheetsParams{
			LoglineID: logline.ID,
		})

		require.NoError(t, err)

		beatsSheets, ok := rawBeatsSheets.(*codegen.GetBeatsSheetsOKApplicationJSON)
		require.True(t, ok)

		require.Len(t, *beatsSheets, 3)
		require.Equal(t, codegen.BeatsSheetPreview{
			ID:        beatsSheet.ID,
			CreatedAt: beatsSheet.CreatedAt,
		}, (*beatsSheets)[0])
	}
}
