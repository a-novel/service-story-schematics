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

func TestListBeatsSheets(t *testing.T) {
	testCases := []struct {
		name string

		fixtures []*dao.BeatsSheetEntity

		data dao.ListBeatsSheetsData

		expect    []*dao.BeatsSheetPreviewEntity
		expectErr error
	}{
		{
			name: "Success",

			fixtures: []*dao.BeatsSheetEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
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
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Content: []models.Beat{
						{
							Key:     "test-beat",
							Title:   "Test Beat",
							Content: "Test Beat Content",
						},
					},
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Content: []models.Beat{
						{
							Key:     "test-beat-2",
							Title:   "Test Beat 2",
							Content: "Test Beat Content 2",
						},
					},
					Lang:      models.LangFR,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.ListBeatsSheetsData{
				LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
			},

			expect: []*dao.BeatsSheetPreviewEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Lang:      models.LangFR,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Limit",

			fixtures: []*dao.BeatsSheetEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
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
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Content: []models.Beat{
						{
							Key:     "test-beat",
							Title:   "Test Beat",
							Content: "Test Beat Content",
						},
					},
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Content: []models.Beat{
						{
							Key:     "test-beat-2",
							Title:   "Test Beat 2",
							Content: "Test Beat Content 2",
						},
					},
					Lang:      models.LangFR,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.ListBeatsSheetsData{
				LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Limit:     2,
			},

			expect: []*dao.BeatsSheetPreviewEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Lang:      models.LangFR,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Offset",

			fixtures: []*dao.BeatsSheetEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
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
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Content: []models.Beat{
						{
							Key:     "test-beat",
							Title:   "Test Beat",
							Content: "Test Beat Content",
						},
					},
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Content: []models.Beat{
						{
							Key:     "test-beat-2",
							Title:   "Test Beat 2",
							Content: "Test Beat Content 2",
						},
					},
					Lang:      models.LangFR,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.ListBeatsSheetsData{
				LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Offset:    1,
			},

			expect: []*dao.BeatsSheetPreviewEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Lang:      models.LangFR,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "SkipOtherLoglines",

			fixtures: []*dao.BeatsSheetEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
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
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					LoglineID: uuid.MustParse("00000000-0000-0000-0001-000000000002"),
					Content: []models.Beat{
						{
							Key:     "test-beat",
							Title:   "Test Beat",
							Content: "Test Beat Content",
						},
					},
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Content: []models.Beat{
						{
							Key:     "test-beat-2",
							Title:   "Test Beat 2",
							Content: "Test Beat Content 2",
						},
					},
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.ListBeatsSheetsData{
				LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
			},

			expect: []*dao.BeatsSheetPreviewEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	repository := dao.NewListBeatsSheetsRepository()

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

			res, err := repository.ListBeatsSheets(tx, testCase.data)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)
		})
	}
}
