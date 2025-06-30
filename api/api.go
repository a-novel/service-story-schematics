package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/ogen-go/ogen/ogenerrors"

	authapi "github.com/a-novel/service-authentication/api"

	"github.com/a-novel/service-story-schematics/api/codegen"
)

var ErrUnauthorized = &codegen.UnexpectedErrorStatusCode{
	StatusCode: http.StatusUnauthorized,
	Response:   codegen.UnexpectedError{Error: "Unauthorized"},
}

var ErrForbidden = &codegen.UnexpectedErrorStatusCode{
	StatusCode: http.StatusForbidden,
	Response:   codegen.UnexpectedError{Error: "Forbidden"},
}

var ErrInternalServerError = &codegen.UnexpectedErrorStatusCode{
	StatusCode: http.StatusInternalServerError,
	Response:   codegen.UnexpectedError{Error: "internal server error"},
}

type API struct {
	codegen.UnimplementedHandler

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
}

func (api *API) NewError(ctx context.Context, err error) *codegen.UnexpectedErrorStatusCode {
	// no-op
	if err == nil {
		return nil
	}

	// Return a different error if authentication failed. Also do not log error (we will still have the API log from
	// the default middleware if needed).
	var securityError *ogenerrors.SecurityError
	if ok := errors.As(err, &securityError); ok {
		switch {
		case errors.Is(err, authapi.ErrAuthentication):
			return ErrUnauthorized
		case errors.Is(err, authapi.ErrPermission):
			return ErrForbidden
		default:
			return ErrUnauthorized
		}
	}

	logger := sentry.NewLogger(ctx)
	logger.Errorf(ctx, "security error: %v", err)

	return ErrInternalServerError
}
