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

func TestSelectBeatsSheet(t *testing.T) {
	testCases := []struct {
		name string

		fixtures []*dao.BeatsSheetEntity

		data uuid.UUID

		expect    *dao.BeatsSheetEntity
		expectErr error
	}{
		{
			name: "Success",

			fixtures: []*dao.BeatsSheetEntity{
				{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					LoglineID:   uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					StoryPlanID: uuid.MustParse("00000000-0000-0000-0002-000000000001"),
					Content: []models.Beat{
						{
							Key:     "test-beat",
							Title:   "Test Beat",
							Content: "Test Beat Content",
						},
						{
							Key:     "test-beat-2",
							Title:   "Test Beat 2",
							Content: "Test Beat Content 2",
						},
					},
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			data: uuid.MustParse("00000000-0000-0000-0000-000000000001"),

			expect: &dao.BeatsSheetEntity{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				LoglineID:   uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				StoryPlanID: uuid.MustParse("00000000-0000-0000-0002-000000000001"),
				Content: []models.Beat{
					{
						Key:     "test-beat",
						Title:   "Test Beat",
						Content: "Test Beat Content",
					},
					{
						Key:     "test-beat-2",
						Title:   "Test Beat 2",
						Content: "Test Beat Content 2",
					},
				},
				Lang:      models.LangEN,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "NotFound",

			data: uuid.MustParse("00000000-0000-0000-0000-000000000001"),

			expectErr: dao.ErrBeatsSheetNotFound,
		},
	}

	repository := dao.NewSelectBeatsSheetRepository()

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

			res, err := repository.SelectBeatsSheet(tx, testCase.data)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)
		})
	}
}
