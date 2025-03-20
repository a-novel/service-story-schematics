package dao_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pgctx "github.com/a-novel-kit/context/pgbun"

	"github.com/a-novel/story-schematics/internal/dao"
)

func TestSelectLoglineBySlug(t *testing.T) {
	testCases := []struct {
		name string

		fixtures []*dao.LoglineEntity

		data dao.SelectLoglineBySlugData

		expect    *dao.LoglineEntity
		expectErr error
	}{
		{
			name: "Success",

			fixtures: []*dao.LoglineEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.SelectLoglineBySlugData{
				Slug:   "test-slug",
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
			},

			expect: &dao.LoglineEntity{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:      "test-slug",
				Name:      "Test Name 2",
				Content:   "Lorem ipsum dolor sit amet 2",
				CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "NotFound/UserID",

			fixtures: []*dao.LoglineEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.SelectLoglineBySlugData{
				Slug:   "test-slug",
				UserID: uuid.MustParse("00000000-0000-0000-0001-000000000002"),
			},

			expectErr: dao.ErrLoglineNotFound,
		},
		{
			name: "NotFound/Slug",

			fixtures: []*dao.LoglineEntity{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			data: dao.SelectLoglineBySlugData{
				Slug:   "test-slug-2",
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
			},

			expectErr: dao.ErrLoglineNotFound,
		},
	}

	repository := dao.NewSelectLoglineBySlugRepository()

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

			res, err := repository.SelectLoglineBySlug(tx, testCase.data)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)
		})
	}
}
