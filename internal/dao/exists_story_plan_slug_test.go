package dao_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/postgres"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/config"
)

func TestExistsStoryPlanSlug(t *testing.T) {
	testCases := []struct {
		name string

		fixtures []*dao.StoryPlanEntity

		data models.Slug

		expect    bool
		expectErr error
	}{
		{
			name: "Success",

			fixtures: []*dao.StoryPlanEntity{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug: "test-slug",

					Name:        "Test Name",
					Description: "Test Description, a lot going on here.",

					Beats: []models.BeatDefinition{
						{
							Name:      "Test Beat",
							Key:       "test-beat",
							KeyPoints: []string{"The key point of the beat, in a single sentence."},
							Purpose:   "The purpose of the beat, in a single sentence.",
						},
						{
							Name:      "Test Beat 2",
							Key:       "test-beat-2",
							KeyPoints: []string{"The key point of the second beat, in a single sentence."},
							Purpose:   "The purpose of the plot second point, in a single sentence.",
						},
					},
					Lang: models.LangEN,

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			data: "test-slug",

			expect: true,
		},
		{
			name: "MultiplePlans",

			fixtures: []*dao.StoryPlanEntity{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug: "test-slug",

					Name:        "Test Name",
					Description: "Test Description, a lot going on here.",

					Beats: []models.BeatDefinition{
						{
							Name:      "Test Beat",
							Key:       "test-beat",
							KeyPoints: []string{"The key point of the beat, in a single sentence."},
							Purpose:   "The purpose of the beat, in a single sentence.",
						},
						{
							Name:      "Test Beat 2",
							Key:       "test-beat-2",
							KeyPoints: []string{"The key point of the second beat, in a single sentence."},
							Purpose:   "The purpose of the plot second point, in a single sentence.",
						},
					},
					Lang: models.LangEN,

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				// Most recent.
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Slug: "test-slug",

					Name:        "Test Name 2",
					Description: "Test Description 2, a lot going on here.",

					Beats: []models.BeatDefinition{
						{
							Name:      "Test Beat 2",
							Key:       "test-beat-2",
							KeyPoints: []string{"The key point of the second beat, in a single sentence."},
							Purpose:   "The purpose of the plot second point, in a single sentence.",
						},
					},
					Lang: models.LangEN,

					CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Slug: "test-slug",

					Name:        "Test Name 3",
					Description: "Test Description 3, a lot going on here.",

					Beats: []models.BeatDefinition{
						{
							Name:      "Test Beat",
							Key:       "test-beat",
							KeyPoints: []string{"The key point of the beat, in a single sentence."},
							Purpose:   "The purpose of the beat, in a single sentence.",
						},
					},
					Lang: models.LangEN,

					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},

			data: "test-slug",

			expect: true,
		},
		{
			name: "NotFound",

			fixtures: []*dao.StoryPlanEntity{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug: "test-slug",

					Name:        "Test Name",
					Description: "Test Description, a lot going on here.",

					Beats: []models.BeatDefinition{
						{
							Name:      "Test Beat",
							Key:       "test-beat",
							KeyPoints: []string{"The key point of the beat, in a single sentence."},
							Purpose:   "The purpose of the beat, in a single sentence.",
						},
						{
							Name:      "Test Beat 2",
							Key:       "test-beat-2",
							KeyPoints: []string{"The key point of the second beat, in a single sentence."},
							Purpose:   "The purpose of the plot second point, in a single sentence.",
						},
					},
					Lang: models.LangEN,

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			data: "fake-slug",
		},
	}

	repository := dao.NewExistsStoryPlanSlugRepository()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			postgres.RunTransactionalTest(t, config.PostgresPresetTest, func(ctx context.Context, t *testing.T) {
				t.Helper()

				db, err := postgres.GetContext(ctx)
				require.NoError(t, err)

				if len(testCase.fixtures) > 0 {
					_, err = db.NewInsert().Model(&testCase.fixtures).Exec(ctx)
					require.NoError(t, err)
				}

				res, err := repository.ExistsStoryPlanSlug(ctx, testCase.data)
				require.ErrorIs(t, err, testCase.expectErr)
				require.Equal(t, testCase.expect, res)
			})
		})
	}
}
