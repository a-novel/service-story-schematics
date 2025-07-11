package services_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/services"
	servicesmocks "github.com/a-novel/service-story-schematics/internal/services/mocks"
	"github.com/a-novel/service-story-schematics/models"
)

func TestExpandLogline(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type expandLoglineData struct {
		resp *models.LoglineIdea
		err  error
	}

	testCases := []struct {
		name string

		request services.ExpandLoglineRequest

		expandLoglineData *expandLoglineData

		expect    *models.LoglineIdea
		expectErr error
	}{
		{
			name: "Success",

			request: services.ExpandLoglineRequest{
				Logline: models.LoglineIdea{
					Name:    "test title",
					Content: "test content",
					Lang:    models.LangEN,
				},
				UserID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},

			expandLoglineData: &expandLoglineData{
				resp: &models.LoglineIdea{
					Name:    "test",
					Content: "test",
					Lang:    models.LangEN,
				},
			},

			expect: &models.LoglineIdea{
				Name:    "test",
				Content: "test",
				Lang:    models.LangEN,
			},
		},
		{
			name: "Error",

			request: services.ExpandLoglineRequest{
				Logline: models.LoglineIdea{
					Name:    "test title",
					Content: "test content",
					Lang:    models.LangEN,
				},
				UserID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},

			expandLoglineData: &expandLoglineData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			source := servicesmocks.NewMockExpandLoglineSource(t)

			if testCase.expandLoglineData != nil {
				source.EXPECT().
					ExpandLogline(mock.Anything, daoai.ExpandLoglineRequest{
						Logline: testCase.request.Logline.Name + "\n\n" + testCase.request.Logline.Content,
						UserID:  testCase.request.UserID.String(),
						Lang:    testCase.request.Logline.Lang,
					}).
					Return(testCase.expandLoglineData.resp, testCase.expandLoglineData.err)
			}

			service := services.NewExpandLoglineService(source)

			resp, err := service.ExpandLogline(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
