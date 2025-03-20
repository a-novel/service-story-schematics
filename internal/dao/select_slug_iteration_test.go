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
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.SelectSlugIterationData{
				Slug: "test-slug",

				Table: "loglines",

				Filter: map[string][]any{
					"user_id = ?": {uuid.MustParse("00000000-0000-0000-1000-000000000001")},
				},

				Order: []string{"created_at DESC"},
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
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				&dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-1",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				&dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-2",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.SelectSlugIterationData{
				Slug: "test-slug",

				Table: "loglines",

				Filter: map[string][]any{
					"user_id = ?": {uuid.MustParse("00000000-0000-0000-1000-000000000001")},
				},

				Order: []string{"created_at DESC"},
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
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				&dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-3",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				&dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-100",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.SelectSlugIterationData{
				Slug: "test-slug",

				Table: "loglines",

				Filter: map[string][]any{
					"user_id = ?": {uuid.MustParse("00000000-0000-0000-1000-000000000001")},
				},

				Order: []string{"created_at DESC"},
			},

			expect:          "test-slug-101",
			expectIteration: 101,
		},
	}

	repository := dao.NewSelectSlugIterationRepository()

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

			res, iter, err := repository.SelectSlugIteration(tx, testCase.data)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expectIteration, iter)
			require.Equal(t, testCase.expect, res)
		})
	}
}
