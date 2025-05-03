package api

import (
	"errors"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/rs/zerolog"

	"github.com/a-novel-kit/context"
	sentryctx "github.com/a-novel-kit/context/sentry"

	"github.com/a-novel/service-story-schematics/api/codegen"
)

type API struct {
	CreateBeatsSheetService CreateBeatsSheetService
	CreateLoglineService    CreateLoglineService
	CreateStoryPlanService  CreateStoryPlanService

	ExpandBeatService    ExpandBeatService
	ExpandLoglineService ExpandLoglineService

	GenerateBeatsSheetService GenerateBeatsSheetService
	GenerateLoglinesService   GenerateLoglinesService

	ListBeatsSheetsService ListBeatsSheetsService
	ListLoglinesService    ListLoglinesService
	ListStoryPlansService  ListStoryPlansService

	RegenerateBeatsService RegenerateBeatsService

	SelectBeatsSheetService SelectBeatsSheetService
	SelectLoglineService    SelectLoglineService
	SelectStoryPlanService  SelectStoryPlanService

	UpdateStoryPlanService UpdateStoryPlanService

	codegen.UnimplementedHandler
}

func (api *API) NewError(ctx context.Context, err error) *codegen.UnexpectedErrorStatusCode {
	// no-op
	if err == nil {
		return nil
	}

	logger := zerolog.Ctx(ctx)

	// Return a different error if authentication failed. Also do not log error (we will still have the API log from
	// the default middleware if needed).
	var securityError *ogenerrors.SecurityError
	if ok := errors.As(err, &securityError); ok {
		logger.Warn().Err(err).Msg("authentication failed")

		return &codegen.UnexpectedErrorStatusCode{
			StatusCode: http.StatusUnauthorized,
			Response:   codegen.UnexpectedError{Error: "Unauthorized"},
		}
	}

	// Unhandled, unexpected error occurred.
	logger.Error().Err(err).Msg("internal error")
	sentryctx.CaptureException(ctx, err)

	return &codegen.UnexpectedErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response:   codegen.UnexpectedError{Error: "internal server error"},
	}
}
