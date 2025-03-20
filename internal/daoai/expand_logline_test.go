package daoai_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/golm"
	groqbinding "github.com/a-novel-kit/golm/bindings/groq"

	"github.com/a-novel/story-schematics/config"
	"github.com/a-novel/story-schematics/internal/daoai"
)

func TestExpandLogline(t *testing.T) {
	testCases := []struct {
		name string

		request daoai.ExpandLoglineRequest
	}{
		{
			name: "Success",

			request: daoai.ExpandLoglineRequest{
				Logline: `The Aurora Initiative

As a team of scientists discover a way to harness the energy of a nearby supernova, they must also contend with the 
implications of altering the course of human history and the emergence of a new, technologically advanced world order.`,
			},
		},
	}

	binding := groqbinding.New(config.Groq.APIKey, config.Groq.Model)

	repository := daoai.NewExpandLoglineRepository()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := golm.WithContext(t.Context(), binding)

			resp, err := repository.ExpandLogline(ctx, testCase.request)
			require.NoError(t, err)

			require.NotNil(t, resp)

			require.NotEmpty(t, resp.Name)
			require.NotEmpty(t, resp.Content)

			CheckAgent(
				t,
				fmt.Sprintf(
					"Does this logline:\n\n%s\n\nExpands this one while retaining its original themes:\n\n%s",
					resp.Content, testCase.request.Logline,
				),
				fmt.Sprintf(
					"The greater AI decreted that this logline:\n\n%s\n\nDoes not expand this one:\n\n%s",
					resp.Content, testCase.request.Logline,
				),
			)
		})
	}
}
