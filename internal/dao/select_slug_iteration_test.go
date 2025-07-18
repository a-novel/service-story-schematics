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

func TestSelectSlugIteration(t *testing.T) {
	testCases := []struct {
		name string

		fixtures []any

		data dao.SelectSlugIterationData

		expect          models.Slug
		expectIteration int
		expectErr       error
	}{
		{
			name: "Success",

			fixtures: []any{
				&dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.SelectSlugIterationData{
				Slug: "test-slug",

				Target: dao.SlugIterationTargetLogline,

				Args: []any{
					uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				},
			},

			expect:          "test-slug-1",
			expectIteration: 1,
		},
		{
			name: "MultipleIterations",

			fixtures: []any{
				&dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				&dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-1",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				&dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-2",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.SelectSlugIterationData{
				Slug: "test-slug",

				Target: dao.SlugIterationTargetLogline,

				Args: []any{
					uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				},
			},

			expect:          "test-slug-3",
			expectIteration: 3,
		},
		{
			name: "BigNumbers",

			fixtures: []any{
				&dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				&dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-3",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				&dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-100",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.SelectSlugIterationData{
				Slug: "test-slug",

				Target: dao.SlugIterationTargetLogline,

				Args: []any{
					uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				},
			},

			expect:          "test-slug-101",
			expectIteration: 101,
		},
	}

	repository := dao.NewSelectSlugIterationRepository()

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

				res, iter, err := repository.SelectSlugIteration(ctx, testCase.data)
				require.ErrorIs(t, err, testCase.expectErr)
				require.Equal(t, testCase.expectIteration, iter)
				require.Equal(t, testCase.expect, res)
			})
		})
	}
}
