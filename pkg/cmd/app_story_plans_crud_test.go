package cmdpkg_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	authmodels "github.com/a-novel/service-authentication/models"

	apimodels "github.com/a-novel/service-story-schematics/models/api"
	"github.com/a-novel/service-story-schematics/pkg"
)

var saveTheCatPartialPlanForm = &apimodels.CreateStoryPlanForm{
	Name: "Save The Cat Simplified",
	Description: `The "Save The Cat" simplified story plan consists of 5 beats that serve as a
blueprint for crafting compelling stories.`,
	Beats: []apimodels.BeatDefinition{
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
	Lang: apimodels.LangEn,
}

func testAppStoryPlansCRUD(ctx context.Context, t *testing.T, appConfig TestConfig) {
	t.Helper()

	security := pkg.NewBearerSource()

	client, err := pkg.NewAPIClient(ctx, fmt.Sprintf("http://localhost:%v/v1", appConfig.API.Port), security)
	require.NoError(t, err)

	storyPlanSlug := "story-plan-crud-integration-test"
	planForm := *saveTheCatPartialPlanForm
	planForm.Slug = apimodels.Slug(storyPlanSlug)

	userLambdaClaims := authmodels.AccessTokenClaims{
		UserID: lo.ToPtr(uuid.New()),
		Roles:  []authmodels.Role{authmodels.RoleUser},
	}
	userSuperAdminClaims := authmodels.AccessTokenClaims{
		UserID: lo.ToPtr(uuid.New()),
		Roles:  []authmodels.Role{authmodels.RoleSuperAdmin},
	}

	userLambdaAccessToken := getAccessToken(t, appConfig, userLambdaClaims)
	userSuperAdminAccessToken := getAccessToken(t, appConfig, userSuperAdminClaims)

	t.Log("CreateStoryPlanNotAllowed")
	{
		security.SetToken(userLambdaAccessToken)

		rawRes, err := client.CreateStoryPlan(t.Context(), &planForm)

		require.NoError(t, err)

		_, ok := rawRes.(*apimodels.ForbiddenError)
		require.True(t, ok, rawRes)
	}

	t.Log("CreateStoryPlan")
	{
		security.SetToken(userSuperAdminAccessToken)

		planRaw, err := client.CreateStoryPlan(t.Context(), &planForm)
		require.NoError(t, err)

		plan, ok := planRaw.(*apimodels.StoryPlan)
		require.True(t, ok, planRaw)
		require.NotEmpty(t, plan.GetID())
		require.Equal(t, planForm.Slug, plan.GetSlug())
		require.Equal(t, planForm.Name, plan.GetName())
		require.Equal(t, planForm.Description, plan.GetDescription())
		require.Equal(t, plan.GetBeats(), planForm.Beats)
	}

	t.Log("CreateStoryPlan/SlugResolution")
	{
		security.SetToken(userSuperAdminAccessToken)

		planRaw, err := client.CreateStoryPlan(t.Context(), &planForm)
		require.NoError(t, err)

		plan, ok := planRaw.(*apimodels.StoryPlan)
		require.True(t, ok, planRaw)
		require.NotEmpty(t, plan.GetID())
		require.Equal(t, planForm.Slug+"-1", plan.GetSlug())
		require.Equal(t, planForm.Name, plan.GetName())
		require.Equal(t, planForm.Description, plan.GetDescription())
		require.Equal(t, plan.GetBeats(), planForm.Beats)
	}

	t.Log("ListStoryPlans")
	{
		security.SetToken(userLambdaAccessToken)

		rawRes, err := client.GetStoryPlans(t.Context(), apimodels.GetStoryPlansParams{})
		require.NoError(t, err)

		res, ok := rawRes.(*apimodels.GetStoryPlansOKApplicationJSON)
		require.True(t, ok, rawRes)

		versions := lo.Filter(*res, func(item apimodels.StoryPlanPreview, _ int) bool {
			return strings.HasPrefix(string(item.GetSlug()), storyPlanSlug)
		})

		require.Len(t, versions, 2)
		require.Equal(t, apimodels.Slug(storyPlanSlug), versions[0].GetSlug())
		require.Equal(t, apimodels.Slug(storyPlanSlug+"-1"), versions[1].GetSlug())
	}

	t.Log("GetStoryPlan")
	{
		security.SetToken(userLambdaAccessToken)

		rawRes, err := client.GetStoryPlan(t.Context(), apimodels.GetStoryPlanParams{
			Slug: apimodels.OptSlug{Value: apimodels.Slug(storyPlanSlug), Set: true},
		})
		require.NoError(t, err)

		res, ok := rawRes.(*apimodels.StoryPlan)
		require.True(t, ok, rawRes)

		require.Equal(t, apimodels.Slug(storyPlanSlug), res.GetSlug())
		require.Equal(t, planForm.Slug, res.GetSlug())
		require.Equal(t, planForm.Name, res.GetName())
		require.Equal(t, planForm.Description, res.GetDescription())
		require.Equal(t, res.GetBeats(), planForm.Beats)
	}

	t.Log("UpdateStoryPlanNotAllowed")
	{
		security.SetToken(userLambdaAccessToken)

		rawRes, err := client.UpdateStoryPlan(t.Context(), &apimodels.UpdateStoryPlanForm{
			Slug:        apimodels.Slug(storyPlanSlug),
			Name:        planForm.Name + " Updated",
			Description: planForm.Description + " Updated",
			Lang:        apimodels.LangEn,
			Beats:       planForm.Beats,
		})

		require.NoError(t, err)

		_, ok := rawRes.(*apimodels.ForbiddenError)
		require.True(t, ok, rawRes)
	}

	t.Log("UpdateStoryPlan")
	{
		security.SetToken(userSuperAdminAccessToken)

		planRaw, err := client.UpdateStoryPlan(t.Context(), &apimodels.UpdateStoryPlanForm{
			Slug:        apimodels.Slug(storyPlanSlug),
			Name:        planForm.Name + " Updated",
			Description: planForm.Description + " Updated",
			Lang:        apimodels.LangEn,
			Beats:       planForm.Beats,
		})
		require.NoError(t, err)

		plan, ok := planRaw.(*apimodels.StoryPlan)
		require.True(t, ok, planRaw)
		require.NotEmpty(t, plan.GetID())
		require.Equal(t, planForm.Slug, plan.GetSlug())
		require.Equal(t, planForm.Name+" Updated", plan.GetName())
		require.Equal(t, planForm.Description+" Updated", plan.GetDescription())
		require.Equal(t, plan.GetBeats(), planForm.Beats)
	}
}
