package dao_test

import (
	"database/sql"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

func TestInsertStoryPlan(t *testing.T) {
	testCases := []struct {
		name string

		fixtures []*dao.StoryPlanEntity

		data dao.InsertStoryPlanData

		expect    *dao.StoryPlanEntity
		expectErr error
	}{
		{
			name: "Success",

			data: dao.InsertStoryPlanData{
				Plan: models.StoryPlan{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug:        "test-slug",
					Name:        "Test Name",
					Description: "Test Description, a lot going on here.",
					Beats: []models.BeatDefinition{
						{
							Name:      "Test Beat",
							Key:       "test-beat",
							KeyPoints: []string{"The key point of the beat, in a single sentence."},
							Purpose:   "The purpose of the beat, in a single sentence.",
							MinScenes: 1,
						},
						{
							Name:      "Test Beat 2",
							Key:       "test-beat-2",
							KeyPoints: []string{"The key point of the second beat, in a single sentence."},
							Purpose:   "The purpose of the plot second point, in a single sentence.",
							MinScenes: 3,
							MaxScenes: 5,
						},
					},
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &dao.StoryPlanEntity{
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
						MinScenes: 1,
					},
					{
						Name:      "Test Beat 2",
						Key:       "test-beat-2",
						KeyPoints: []string{"The key point of the second beat, in a single sentence."},
						Purpose:   "The purpose of the plot second point, in a single sentence.",
						MinScenes: 3,
						MaxScenes: 5,
					},
				},
				Lang: models.LangEN,

				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Conflict",

			data: dao.InsertStoryPlanData{
				Plan: models.StoryPlan{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug:        "test-slug",
					Name:        "Test Name",
					Description: "Test Description, a lot going on here.",
					Beats: []models.BeatDefinition{
						{
							Name:      "Test Beat",
							Key:       "test-beat",
							KeyPoints: []string{"The key point of the beat, in a single sentence."},
							Purpose:   "The purpose of the beat, in a single sentence.",
							MinScenes: 1,
						},
						{
							Name:      "Test Beat 2",
							Key:       "test-beat-2",
							KeyPoints: []string{"The key point of the beat, in a single sentence."},
							Purpose:   "The purpose of the plot second point, in a single sentence.",
							MinScenes: 3,
							MaxScenes: 5,
						},
					},
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			fixtures: []*dao.StoryPlanEntity{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Slug: "test-slug",

					Name:        "Test Name 2",
					Description: "Test Description 2, a lot going on here.",

					Beats: []models.BeatDefinition{
						{
							Name:      "Test Beat 3",
							Key:       "test-beat 3",
							KeyPoints: []string{"The key point of the beat, in a single sentence."},
							Purpose:   "The purpose of the beat, in a single sentence.",
							MinScenes: 1,
						},
					},
					Lang: models.LangEN,

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expectErr: dao.ErrStoryPlanAlreadyExists,
		},
	}

	repository := dao.NewInsertStoryPlanRepository()

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

			res, err := repository.InsertStoryPlan(tx, testCase.data)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)
		})
	}
}
