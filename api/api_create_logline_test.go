package api_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	authapi "github.com/a-novel/service-authentication/api"
	authmodels "github.com/a-novel/service-authentication/models"

	"github.com/a-novel/service-story-schematics/api"
	"github.com/a-novel/service-story-schematics/api/codegen"
	apimocks "github.com/a-novel/service-story-schematics/api/mocks"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

func TestCreateLogline(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type createLoglineData struct {
		resp *models.Logline
		err  error
	}

	testCases := []struct {
		name string

		form *codegen.CreateLoglineForm

		createLoglineData *createLoglineData

		expect    codegen.CreateLoglineRes
		expectErr error
	}{
		{
			name: "Success",

			form: &codegen.CreateLoglineForm{
				Slug:    "slug",
				Name:    "name",
				Content: "content",
				Lang:    codegen.LangEn,
			},

			createLoglineData: &createLoglineData{
				resp: &models.Logline{
					ID:        uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug:      "slug",
					Name:      "name",
					Content:   "content",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &codegen.Logline{
				ID:        codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				UserID:    codegen.UserID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				Slug:      "slug",
				Name:      "name",
				Content:   "content",
				Lang:      codegen.LangEn,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Error",

			form: &codegen.CreateLoglineForm{
				Slug:    "slug",
				Name:    "name",
				Content: "content",
				Lang:    codegen.LangEn,
			},

			createLoglineData: &createLoglineData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			source := apimocks.NewMockCreateLoglineService(t)

			ctx := context.WithValue(t.Context(), authapi.ClaimsAPIKey{}, &authmodels.AccessTokenClaims{
				UserID: lo.ToPtr(uuid.MustParse("00000000-1000-0000-0000-000000000001")),
			})

			if testCase.createLoglineData != nil {
				source.EXPECT().
					CreateLogline(mock.Anything, services.CreateLoglineRequest{
						UserID:  uuid.MustParse("00000000-1000-0000-0000-000000000001"),
						Slug:    models.Slug(testCase.form.GetSlug()),
						Name:    testCase.form.GetName(),
						Content: testCase.form.GetContent(),
						Lang:    models.Lang(testCase.form.GetLang()),
					}).
					Return(testCase.createLoglineData.resp, testCase.createLoglineData.err)
			}

			handler := api.API{CreateLoglineService: source}

			res, err := handler.CreateLogline(ctx, testCase.form)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
