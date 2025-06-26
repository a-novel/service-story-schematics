package dao_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
)

func TestListStoryPlans(t *testing.T) {
	testCases := []struct {
		name string

		fixtures []*dao.StoryPlanEntity

		data dao.ListStoryPlansData

		expect    []*dao.StoryPlanPreviewEntity
		expectErr error
	}{
		{
			name: "Success",

			fixtures: []*dao.StoryPlanEntity{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug: "test-slug-1",

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
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Slug: "test-slug-2",

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
					Slug: "test-slug-3",

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
					Lang: models.LangFR,

					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.ListStoryPlansData{},

			expect: []*dao.StoryPlanPreviewEntity{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Slug: "test-slug-2",

					Name:        "Test Name 2",
					Description: "Test Description 2, a lot going on here.",
					Lang:        models.LangEN,

					CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Slug: "test-slug-3",

					Name:        "Test Name 3",
					Description: "Test Description 3, a lot going on here.",
					Lang:        models.LangFR,

					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug: "test-slug-1",

					Name:        "Test Name",
					Description: "Test Description, a lot going on here.",
					Lang:        models.LangEN,

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Limit",

			fixtures: []*dao.StoryPlanEntity{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug: "test-slug-1",

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
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Slug: "test-slug-2",

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
					Slug: "test-slug-3",

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
					Lang: models.LangFR,

					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.ListStoryPlansData{
				Limit: 2,
			},

			expect: []*dao.StoryPlanPreviewEntity{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Slug: "test-slug-2",

					Name:        "Test Name 2",
					Description: "Test Description 2, a lot going on here.",
					Lang:        models.LangEN,

					CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Slug: "test-slug-3",

					Name:        "Test Name 3",
					Description: "Test Description 3, a lot going on here.",
					Lang:        models.LangFR,

					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Offset",

			fixtures: []*dao.StoryPlanEntity{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug: "test-slug-1",

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
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Slug: "test-slug-2",

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
					Slug: "test-slug-3",

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
					Lang: models.LangFR,

					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.ListStoryPlansData{
				Offset: 1,
			},

			expect: []*dao.StoryPlanPreviewEntity{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Slug: "test-slug-3",

					Name:        "Test Name 3",
					Description: "Test Description 3, a lot going on here.",
					Lang:        models.LangFR,

					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug: "test-slug-1",

					Name:        "Test Name",
					Description: "Test Description, a lot going on here.",
					Lang:        models.LangEN,

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "NoDuplicates",

			fixtures: []*dao.StoryPlanEntity{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug: "test-slug-1",

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
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Slug: "test-slug-2",

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
					Slug: "test-slug-3",

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
				// Newer version of test-slug-2.
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Slug: "test-slug-2",

					Name:        "Test New Name 2",
					Description: "Test New Description 2, a lot going on here.",

					Beats: []models.BeatDefinition{
						{
							Name:      "Test Beat 2",
							Key:       "test-beat-2",
							KeyPoints: []string{"The key point of the second beat, in a single sentence."},
							Purpose:   "The purpose of the plot second point, in a single sentence.",
						},
						{
							Name:      "Test Beat 4",
							Key:       "test-beat-4",
							KeyPoints: []string{"The key point of the fourth beat, in a single sentence."},
							Purpose:   "The purpose of the plot fourth point, in a single sentence.",
						},
					},
					Lang: models.LangEN,

					CreatedAt: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.ListStoryPlansData{},

			expect: []*dao.StoryPlanPreviewEntity{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000004"),
					Slug: "test-slug-2",

					Name:        "Test New Name 2",
					Description: "Test New Description 2, a lot going on here.",
					Lang:        models.LangEN,

					CreatedAt: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Slug: "test-slug-3",

					Name:        "Test Name 3",
					Description: "Test Description 3, a lot going on here.",
					Lang:        models.LangEN,

					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug: "test-slug-1",

					Name:        "Test Name",
					Description: "Test Description, a lot going on here.",
					Lang:        models.LangEN,

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	repository := dao.NewListStoryPlansRepository()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tx, commit, err := lib.PostgresContextTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
			require.NoError(t, err)

			t.Cleanup(func() {
				_ = commit(false)
			})

			db, err := lib.PostgresContext(tx)
			require.NoError(t, err)

			if len(testCase.fixtures) > 0 {
				_, err = db.NewInsert().Model(&testCase.fixtures).Exec(tx)
				require.NoError(t, err)
			}

			res, err := repository.ListStoryPlans(tx, testCase.data)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)
		})
	}
}
