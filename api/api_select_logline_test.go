package api_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	authapi "github.com/a-novel/authentication/api"
	authmodels "github.com/a-novel/authentication/models"

	"github.com/a-novel-kit/context"

	"github.com/a-novel/story-schematics/api"
	"github.com/a-novel/story-schematics/api/codegen"
	apimocks "github.com/a-novel/story-schematics/api/mocks"
	"github.com/a-novel/story-schematics/internal/dao"
	"github.com/a-novel/story-schematics/internal/services"
	"github.com/a-novel/story-schematics/models"
)

func TestSelectLogline(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type selectLoglineData struct {
		resp *models.Logline
		err  error
	}

	testCases := []struct {
		name string

		params codegen.GetLoglineParams

		selectLoglineData *selectLoglineData

		expect    codegen.GetLoglineRes
		expectErr error
	}{
		{
			name: "Success/ID",

			params: codegen.GetLoglineParams{
				ID: codegen.OptLoglineID{
					Value: codegen.LoglineID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					Set:   true,
				},
				Slug: codegen.OptSlug{},
			},

			selectLoglineData: &selectLoglineData{
				resp: &models.Logline{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &codegen.Logline{
				ID:        codegen.LoglineID(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
				UserID:    codegen.UserID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				Slug:      "test-slug",
				Name:      "Test Name 2",
				Content:   "Lorem ipsum dolor sit amet 2",
				CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Success/Slug",

			params: codegen.GetLoglineParams{
				ID: codegen.OptLoglineID{},
				Slug: codegen.OptSlug{
					Value: "test-slug",
					Set:   true,
				},
			},

			selectLoglineData: &selectLoglineData{
				resp: &models.Logline{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &codegen.Logline{
				ID:        codegen.LoglineID(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
				UserID:    codegen.UserID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				Slug:      "test-slug",
				Name:      "Test Name 2",
				Content:   "Lorem ipsum dolor sit amet 2",
				CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "LoglineNotFound",

			params: codegen.GetLoglineParams{
				ID: codegen.OptLoglineID{
					Value: codegen.LoglineID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					Set:   true,
				},
				Slug: codegen.OptSlug{},
			},

			selectLoglineData: &selectLoglineData{
				err: dao.ErrLoglineNotFound,
			},

			expect: &codegen.NotFoundError{Error: dao.ErrLoglineNotFound.Error()},
		},
		{
			name: "Error",

			params: codegen.GetLoglineParams{
				ID: codegen.OptLoglineID{
					Value: codegen.LoglineID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					Set:   true,
				},
				Slug: codegen.OptSlug{},
			},

			selectLoglineData: &selectLoglineData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			source := apimocks.NewMockSelectLoglineService(t)

			ctx := context.WithValue(t.Context(), authapi.ClaimsAPIKey{}, &authmodels.AccessTokenClaims{
				UserID: lo.ToPtr(uuid.MustParse("00000000-1000-0000-0000-000000000001")),
			})

			if testCase.selectLoglineData != nil {
				source.EXPECT().
					SelectLogline(ctx, services.SelectLoglineRequest{
						UserID: uuid.MustParse("00000000-1000-0000-0000-000000000001"),
						Slug: lo.Ternary(
							testCase.params.Slug.IsSet(),
							lo.ToPtr(models.Slug(testCase.params.Slug.Value)),
							nil,
						),
						ID: lo.Ternary(
							testCase.params.ID.IsSet(),
							lo.ToPtr(uuid.UUID(testCase.params.ID.Value)),
							nil,
						),
					}).
					Return(testCase.selectLoglineData.resp, testCase.selectLoglineData.err)
			}

			handler := api.API{SelectLoglineService: source}

			res, err := handler.GetLogline(ctx, testCase.params)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
