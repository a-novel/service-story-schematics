package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	"github.com/a-novel/golib/otel"
	authmodels "github.com/a-novel/service-authentication/models"
	jkApiModels "github.com/a-novel/service-json-keys/models/api"

	"github.com/a-novel/service-story-schematics/models/api"
)

var ErrUnauthorized = &apimodels.UnexpectedErrorStatusCode{
	StatusCode: http.StatusUnauthorized,
	Response:   apimodels.UnexpectedError{Error: "Unauthorized"},
}

var ErrForbidden = &apimodels.UnexpectedErrorStatusCode{
	StatusCode: http.StatusForbidden,
	Response:   apimodels.UnexpectedError{Error: "Forbidden"},
}

var ErrInternalServerError = &apimodels.UnexpectedErrorStatusCode{
	StatusCode: http.StatusInternalServerError,
	Response:   apimodels.UnexpectedError{Error: "internal server error"},
}

type API struct {
	apimodels.UnimplementedHandler

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

	JKClient *jkApiModels.Client
}

func (api *API) NewError(ctx context.Context, err error) *apimodels.UnexpectedErrorStatusCode {
	// no-op
	if err == nil {
		return nil
	}

	logger := otel.Logger()

	// Return a different error if authentication failed. Also do not log error (we will still have the API log from
	// the default middleware if needed).
	var securityError *ogenerrors.SecurityError
	if ok := errors.As(err, &securityError); ok {
		logger.ErrorContext(ctx, fmt.Sprintf("security error: %v", err))

		switch {
		case errors.Is(err, authmodels.ErrUnauthorized):
			return ErrUnauthorized
		case errors.Is(err, authmodels.ErrForbidden):
			return ErrForbidden
		default:
			return ErrUnauthorized
		}
	}

	logger.ErrorContext(ctx, fmt.Sprintf("internal server error: %v", err))

	return ErrInternalServerError
}
