package api_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	authModels "github.com/a-novel/service-authentication/models"
	authPkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/api"
	"github.com/a-novel/service-story-schematics/internal/api/codegen"
	apimocks "github.com/a-novel/service-story-schematics/internal/api/mocks"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

func TestGenerateLoglines(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type generateLoglinesData struct {
		loglines []models.LoglineIdea
		err      error
	}

	testCases := []struct {
		name string

		form *codegen.GenerateLoglinesForm

		generateLoglinesData *generateLoglinesData

		expect    codegen.GenerateLoglinesRes
		expectErr error
	}{
		{
			name: "Success",

			form: &codegen.GenerateLoglinesForm{
				Count: 10,
				Theme: "theme",
				Lang:  codegen.LangEn,
			},

			generateLoglinesData: &generateLoglinesData{
				loglines: []models.LoglineIdea{
					{
						Name:    "Logline 1",
						Content: "Logline 1 content",
						Lang:    models.LangEN,
					},
					{
						Name:    "Logline 2",
						Content: "Logline 2 content",
						Lang:    models.LangEN,
					},
				},
			},

			expect: &codegen.GenerateLoglinesOKApplicationJSON{
				{
					Name:    "Logline 1",
					Content: "Logline 1 content",
					Lang:    codegen.LangEn,
				},
				{
					Name:    "Logline 2",
					Content: "Logline 2 content",
					Lang:    codegen.LangEn,
				},
			},
		},
		{
			name: "Error",

			form: &codegen.GenerateLoglinesForm{
				Count: 10,
				Theme: "theme",
				Lang:  codegen.LangEn,
			},

			generateLoglinesData: &generateLoglinesData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			source := apimocks.NewMockGenerateLoglinesService(t)

			ctx := context.WithValue(t.Context(), authPkg.ClaimsContextKey{}, &authModels.AccessTokenClaims{
				UserID: lo.ToPtr(uuid.MustParse("00000000-1000-0000-0000-000000000001")),
			})

			if testCase.generateLoglinesData != nil {
				source.EXPECT().
					GenerateLoglines(mock.Anything, services.GenerateLoglinesRequest{
						Count:  testCase.form.GetCount(),
						Theme:  testCase.form.GetTheme(),
						UserID: uuid.MustParse("00000000-1000-0000-0000-000000000001"),
						Lang:   models.Lang(testCase.form.GetLang()),
					}).
					Return(testCase.generateLoglinesData.loglines, testCase.generateLoglinesData.err)
			}

			handler := api.API{GenerateLoglinesService: source}

			res, err := handler.GenerateLoglines(ctx, testCase.form)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
