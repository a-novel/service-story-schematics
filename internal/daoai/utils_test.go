package daoai_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/packages/param"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/daoai/testdata"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/config"
)

const TestUser = "agora-story-schematics-test"

// Remove punctuation and other extra characters from the check agent response.
var checkAgentFormatter = regexp.MustCompile("[^a-zA-Z0-9 ]")

// CheckAgent validates the response of an OpenAI agent to a given prompt.
func CheckAgent(t *testing.T, prompt, message string) {
	t.Helper()

	chatCompletion, err := config.OpenAIPresetDefault.Client().
		Chat.Completions.
		New(t.Context(), openai.ChatCompletionNewParams{
			Model: config.OpenAIPresetDefault.Model,
			User:  param.NewOpt(TestUser),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(testdata.UtilsPrompt.CheckAgent.System),
				openai.UserMessage(prompt),
			},
		})

	require.NoError(t, err)

	// Prevent flakiness issues because LLMS just dont give a fuck.
	formattedResponse := chatCompletion.Choices[0].Message.Content
	formattedResponse = strings.TrimSpace(formattedResponse)
	formattedResponse = strings.ToUpper(formattedResponse)
	formattedResponse = checkAgentFormatter.ReplaceAllString(formattedResponse, "")
	require.Equal(t, testdata.UtilsPrompt.CheckAgent.Expect, formattedResponse, message)
}

func CheckLang(t *testing.T, lang models.Lang, rawAnswer string) {
	t.Helper()

	CheckAgent(
		t,
		fmt.Sprintf(testdata.Langs[lang], rawAnswer),
		fmt.Sprintf("answer is not in expected language %s\n%s", lang, rawAnswer),
	)
}
