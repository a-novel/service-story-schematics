package services_test

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/services"
	servicesmocks "github.com/a-novel/service-story-schematics/internal/services/mocks"
	"github.com/a-novel/service-story-schematics/models"
)

func TestGenerateLoglines(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type generateLoglinesData struct {
		resp []models.LoglineIdea
		err  error
	}

	testCases := []struct {
		name string

		request services.GenerateLoglinesRequest

		generateLoglinesData *generateLoglinesData

		expect    []models.LoglineIdea
		expectErr error
	}{
		{
			name: "Success",

			request: services.GenerateLoglinesRequest{
				Count:  5,
				Theme:  "test-theme",
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Lang:   models.LangEN,
			},

			generateLoglinesData: &generateLoglinesData{
				resp: []models.LoglineIdea{
					{
						Name:    "Logline 1",
						Content: "Content 1",
						Lang:    models.LangEN,
					},
					{
						Name:    "Logline 2",
						Content: "Content 2",
						Lang:    models.LangEN,
					},
				},
			},

			expect: []models.LoglineIdea{
				{
					Name:    "Logline 1",
					Content: "Content 1",
					Lang:    models.LangEN,
				},
				{
					Name:    "Logline 2",
					Content: "Content 2",
					Lang:    models.LangEN,
				},
			},
		},
		{
			name: "Error",

			request: services.GenerateLoglinesRequest{
				Count:  5,
				Theme:  "test-theme",
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Lang:   models.LangEN,
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

			ctx := t.Context()

			source := servicesmocks.NewMockGenerateLoglinesSource(t)

			if testCase.generateLoglinesData != nil {
				source.EXPECT().
					GenerateLoglines(mock.Anything, daoai.GenerateLoglinesRequest{
						Count:  testCase.request.Count,
						Theme:  testCase.request.Theme,
						UserID: testCase.request.UserID.String(),
						Lang:   testCase.request.Lang,
					}).
					Return(testCase.generateLoglinesData.resp, testCase.generateLoglinesData.err)
			}

			service := services.NewGenerateLoglinesService(source)

			resp, err := service.GenerateLoglines(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
