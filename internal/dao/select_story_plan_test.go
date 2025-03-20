package dao_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pgctx "github.com/a-novel-kit/context/pgbun"

	"github.com/a-novel/story-schematics/internal/dao"
	"github.com/a-novel/story-schematics/models"
)

func TestSelectStoryPlan(t *testing.T) {
	testCases := []struct {
		name string

		fixtures []*dao.StoryPlanEntity

		data uuid.UUID

		expect    *dao.StoryPlanEntity
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

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			data: uuid.MustParse("00000000-0000-0000-0000-000000000001"),

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
					},
					{
						Name:      "Test Beat 2",
						Key:       "test-beat-2",
						KeyPoints: []string{"The key point of the second beat, in a single sentence."},
						Purpose:   "The purpose of the plot second point, in a single sentence.",
					},
				},

				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
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

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			data: uuid.MustParse("00000000-0000-0000-0000-000000000002"),

			expectErr: dao.ErrStoryPlanNotFound,
		},
	}

	repository := dao.NewSelectStoryPlanRepository()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tx, commit, err := pgctx.NewContextTX(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
			require.NoError(t, err)

			t.Cleanup(func() {
				_ = commit(false)
			})

			db, err := pgctx.Context(tx)
			require.NoError(t, err)

			if len(testCase.fixtures) > 0 {
				_, err = db.NewInsert().Model(&testCase.fixtures).Exec(tx)
				require.NoError(t, err)
			}

			res, err := repository.SelectStoryPlan(tx, testCase.data)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)
		})
	}
}
