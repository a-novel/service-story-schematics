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

func TestInserLogline(t *testing.T) {
	testCases := []struct {
		name string

		fixtures []*dao.LoglineEntity

		data dao.InsertLoglineData

		expect    *dao.LoglineEntity
		expectErr error
	}{
		{
			name: "Success",

			data: dao.InsertLoglineData{
				ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:  uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:    "test-slug",
				Name:    "Test Name",
				Content: "Lorem ipsum dolor sit amet",
				Lang:    models.LangEN,
				Now:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},

			expect: &dao.LoglineEntity{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:      "test-slug",
				Name:      "Test Name",
				Content:   "Lorem ipsum dolor sit amet",
				Lang:      models.LangEN,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "AlreadyExists",

			fixtures: []*dao.LoglineEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.InsertLoglineData{
				ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:  uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:    "test-slug",
				Name:    "Test Name",
				Content: "Lorem ipsum dolor sit amet",
				Lang:    models.LangEN,
				Now:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},

			expectErr: dao.ErrLoglineAlreadyExists,
		},
		{
			name: "SlugTakenDifferentUser",

			fixtures: []*dao.LoglineEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-0001-000000000002"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.InsertLoglineData{
				ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:  uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:    "test-slug",
				Name:    "Test Name",
				Content: "Lorem ipsum dolor sit amet",
				Lang:    models.LangEN,
				Now:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},

			expect: &dao.LoglineEntity{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:      "test-slug",
				Name:      "Test Name",
				Content:   "Lorem ipsum dolor sit amet",
				Lang:      models.LangEN,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "SameUserDifferentSlug",

			fixtures: []*dao.LoglineEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-2",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.InsertLoglineData{
				ID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:  uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:    "test-slug",
				Name:    "Test Name",
				Content: "Lorem ipsum dolor sit amet",
				Lang:    models.LangEN,
				Now:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},

			expect: &dao.LoglineEntity{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:      "test-slug",
				Name:      "Test Name",
				Content:   "Lorem ipsum dolor sit amet",
				Lang:      models.LangEN,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	repository := dao.NewInsertLoglineRepository()

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

				res, err := repository.InsertLogline(ctx, testCase.data)
				require.ErrorIs(t, err, testCase.expectErr)
				require.Equal(t, testCase.expect, res)
			})
		})
	}
}
