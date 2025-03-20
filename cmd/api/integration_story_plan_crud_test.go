package main

import (
	"crypto/rand"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/authentication/api/apiclient/testapiclient"
	authmodels "github.com/a-novel/authentication/models"

	"github.com/a-novel/story-schematics/api/codegen"
)

var saveTheCatPartialPlanForm = &codegen.CreateStoryPlanForm{
	Name: "Save The Cat Simplified",
	Description: `The "Save The Cat" simplified story plan consists of 5 beats that serve as a
blueprint for crafting compelling stories.`,
	Beats: []codegen.BeatDefinition{
		{
			Name:      "Opening Image",
			Key:       "openingImage",
			KeyPoints: []string{"Establish the protagonist's world before the journey begins."},
			Purpose: "Sets the tone, mood, and stakes; offers a visual representation " +
				"of the starting point.",
			MinScenes: 1,
			MaxScenes: 1,
		},
		{
			Name:      "Theme Stated",
			Key:       "themeStated",
			KeyPoints: []string{"Introduce the story's central theme or moral."},
			Purpose:   "Often delivered through dialogue; foreshadows the protagonist transformation.",
			MinScenes: 1,
		},
		{
			Name: "Set-Up",
			Key:  "setup",
			KeyPoints: []string{
				"Introduce the main characters.",
				"Showcase the protagonist's flaws or challenges.",
				"Establish the stakes and the world they live in.",
			},
			Purpose:   "Builds empathy and grounds the audience in the story.",
			MinScenes: 3,
			MaxScenes: 5,
		},
		{
			Name:      "Catalyst",
			Key:       "catalyst",
			KeyPoints: []string{"An event that disrupts the status quo."},
			Purpose:   "Propels the protagonist into the main conflict.",
			MaxScenes: 1,
		},
		{
			Name: "Debate",
			Key:  "debate",
			KeyPoints: []string{
				"The protagonist grapples with the decision to embark on the journey.",
				"Highlights internal conflicts and fears.",
			},
			Purpose:   "Adds depth to the character and heightens tension.",
			MinScenes: 2,
			MaxScenes: 3,
		},
	},
}

func TestStoryPlanCRUD(t *testing.T) {
	client, securityClient, err := getServerClient()
	require.NoError(t, err)

	userLambda := rand.Text()
	userSuperAdmin := rand.Text()

	testapiclient.AddPool(userLambda, &authmodels.AccessTokenClaims{
		UserID: lo.ToPtr(uuid.New()),
		Roles:  []authmodels.Role{authmodels.RoleUser},
	})
	testapiclient.AddPool(userSuperAdmin, &authmodels.AccessTokenClaims{
		UserID: lo.ToPtr(uuid.New()),
		Roles:  []authmodels.Role{authmodels.RoleSuperAdmin},
	})

	storyPlanSlug := "story-plan-crud-integration-test"
	planForm := *saveTheCatPartialPlanForm
	planForm.Slug = codegen.Slug(storyPlanSlug)

	t.Log("CreateStoryPlanNotAllowed")
	{
		securityClient.SetToken(userLambda)

		_, err = client.CreateStoryPlan(t.Context(), &planForm)
		require.Error(t, err)
	}

	t.Log("CreateStoryPlan")
	{
		securityClient.SetToken(userSuperAdmin)

		planRaw, err := client.CreateStoryPlan(t.Context(), &planForm)
		require.NoError(t, err)

		plan, ok := planRaw.(*codegen.StoryPlan)
		require.True(t, ok)
		require.NotEmpty(t, plan.GetID())
		require.Equal(t, planForm.Slug, plan.GetSlug())
		require.Equal(t, planForm.Name, plan.GetName())
		require.Equal(t, planForm.Description, plan.GetDescription())
		require.Equal(t, plan.GetBeats(), planForm.Beats)
	}

	t.Log("CreateStoryPlan/SlugResolution")
	{
		securityClient.SetToken(userSuperAdmin)

		planRaw, err := client.CreateStoryPlan(t.Context(), &planForm)
		require.NoError(t, err)

		plan, ok := planRaw.(*codegen.StoryPlan)
		require.True(t, ok)
		require.NotEmpty(t, plan.GetID())
		require.Equal(t, planForm.Slug+"-1", plan.GetSlug())
		require.Equal(t, planForm.Name, plan.GetName())
		require.Equal(t, planForm.Description, plan.GetDescription())
		require.Equal(t, plan.GetBeats(), planForm.Beats)
	}

	t.Log("ListStoryPlans")
	{
		securityClient.SetToken(userLambda)

		rawRes, err := client.GetStoryPlans(t.Context(), codegen.GetStoryPlansParams{})
		require.NoError(t, err)

		res, ok := rawRes.(*codegen.GetStoryPlansOKApplicationJSON)
		require.True(t, ok)

		versions := lo.Filter(*res, func(item codegen.StoryPlanPreview, _ int) bool {
			return strings.HasPrefix(string(item.GetSlug()), storyPlanSlug)
		})

		require.Len(t, versions, 2)
		require.Equal(t, codegen.Slug(storyPlanSlug), versions[0].GetSlug())
		require.Equal(t, codegen.Slug(storyPlanSlug+"-1"), versions[1].GetSlug())
	}

	t.Log("GetStoryPlan")
	{
		securityClient.SetToken(userLambda)

		rawRes, err := client.GetStoryPlan(t.Context(), codegen.GetStoryPlanParams{
			Slug: codegen.OptSlug{Value: codegen.Slug(storyPlanSlug), Set: true},
		})
		require.NoError(t, err)

		res, ok := rawRes.(*codegen.StoryPlan)
		require.True(t, ok)

		require.Equal(t, codegen.Slug(storyPlanSlug), res.GetSlug())
		require.Equal(t, planForm.Slug, res.GetSlug())
		require.Equal(t, planForm.Name, res.GetName())
		require.Equal(t, planForm.Description, res.GetDescription())
		require.Equal(t, res.GetBeats(), planForm.Beats)
	}

	t.Log("UpdateStoryPlanNotAllowed")
	{
		securityClient.SetToken(userLambda)

		_, err = client.UpdateStoryPlan(t.Context(), &codegen.UpdateStoryPlanForm{
			Slug:        codegen.Slug(storyPlanSlug),
			Name:        planForm.Name + " Updated",
			Description: planForm.Description + " Updated",
			Beats:       planForm.Beats,
		})
		require.Error(t, err)
	}

	t.Log("UpdateStoryPlan")
	{
		securityClient.SetToken(userSuperAdmin)

		planRaw, err := client.UpdateStoryPlan(t.Context(), &codegen.UpdateStoryPlanForm{
			Slug:        codegen.Slug(storyPlanSlug),
			Name:        planForm.Name + " Updated",
			Description: planForm.Description + " Updated",
			Beats:       planForm.Beats,
		})
		require.NoError(t, err)

		plan, ok := planRaw.(*codegen.StoryPlan)
		require.True(t, ok)
		require.NotEmpty(t, plan.GetID())
		require.Equal(t, planForm.Slug, plan.GetSlug())
		require.Equal(t, planForm.Name+" Updated", plan.GetName())
		require.Equal(t, planForm.Description+" Updated", plan.GetDescription())
		require.Equal(t, plan.GetBeats(), planForm.Beats)
	}
}
