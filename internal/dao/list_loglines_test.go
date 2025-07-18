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

func TestListLoglines(t *testing.T) {
	testCases := []struct {
		name string

		fixtures []*dao.LoglineEntity

		data dao.ListLoglinesData

		expect    []*dao.LoglinePreviewEntity
		expectErr error
	}{
		{
			name: "Success",

			fixtures: []*dao.LoglineEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-2",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-3",
					Name:      "Test Name 3",
					Content:   "Lorem ipsum dolor sit amet 3",
					Lang:      models.LangFR,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.ListLoglinesData{
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
			},

			expect: []*dao.LoglinePreviewEntity{
				{
					Slug:      "test-slug",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					Slug:      "test-slug-3",
					Name:      "Test Name 3",
					Content:   "Lorem ipsum dolor sit amet 3",
					Lang:      models.LangFR,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					Slug:      "test-slug-2",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Limit",

			fixtures: []*dao.LoglineEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-2",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-3",
					Name:      "Test Name 3",
					Content:   "Lorem ipsum dolor sit amet 3",
					Lang:      models.LangFR,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.ListLoglinesData{
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Limit:  2,
			},

			expect: []*dao.LoglinePreviewEntity{
				{
					Slug:      "test-slug",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					Slug:      "test-slug-3",
					Name:      "Test Name 3",
					Content:   "Lorem ipsum dolor sit amet 3",
					Lang:      models.LangFR,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Offset",

			fixtures: []*dao.LoglineEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-2",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-3",
					Name:      "Test Name 3",
					Content:   "Lorem ipsum dolor sit amet 3",
					Lang:      models.LangFR,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.ListLoglinesData{
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Offset: 1,
			},

			expect: []*dao.LoglinePreviewEntity{
				{
					Slug:      "test-slug-3",
					Name:      "Test Name 3",
					Content:   "Lorem ipsum dolor sit amet 3",
					Lang:      models.LangFR,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					Slug:      "test-slug-2",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "SkipOtherUsers",

			fixtures: []*dao.LoglineEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-0001-000000000002"),
					Slug:      "test-slug-2",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-3",
					Name:      "Test Name 3",
					Content:   "Lorem ipsum dolor sit amet 3",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.ListLoglinesData{
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
			},

			expect: []*dao.LoglinePreviewEntity{
				{
					Slug:      "test-slug",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					Slug:      "test-slug-3",
					Name:      "Test Name 3",
					Content:   "Lorem ipsum dolor sit amet 3",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	repository := dao.NewListLoglinesRepository()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Run(testCase.name, func(t *testing.T) {
				postgres.RunTransactionalTest(t, config.PostgresPresetTest, func(ctx context.Context, t *testing.T) {
					t.Helper()

					db, err := postgres.GetContext(ctx)
					require.NoError(t, err)

					if len(testCase.fixtures) > 0 {
						_, err = db.NewInsert().Model(&testCase.fixtures).Exec(ctx)
						require.NoError(t, err)
					}

					res, err := repository.ListLoglines(ctx, testCase.data)
					require.ErrorIs(t, err, testCase.expectErr)
					require.Equal(t, testCase.expect, res)
				})
			})
		})
	}
}
