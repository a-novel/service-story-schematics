package api

import (
	"context"
	"errors"

	authapi "github.com/a-novel/service-authentication/api"
	authcodegen "github.com/a-novel/service-authentication/api/codegen"
	"github.com/a-novel/service-authentication/models"

	"github.com/a-novel/service-story-schematics/api/codegen"
)

var ErrStorySchematicsSecurityHandler = errors.New("story schematics security handler")

func NewErrStorySchematicsSecurityHandler(err error) error {
	return errors.Join(err, ErrStorySchematicsSecurityHandler)
}

type SecurityHandler struct {
	authHandler *authapi.SecurityHandler
}

func NewSecurity(
	granted models.PermissionsConfig,
	authService authapi.SecurityHandlerService,
) (*SecurityHandler, error) {
	authHandler, err := authapi.NewSecurity(granted, authService)
	if err != nil {
		return nil, NewErrStorySchematicsSecurityHandler(err)
	}

	return &SecurityHandler{
		authHandler: authHandler,
	}, nil
}

func (security *SecurityHandler) HandleBearerAuth(
	ctx context.Context, operationName codegen.OperationName, auth codegen.BearerAuth,
) (context.Context, error) {
	ctx, err := security.authHandler.HandleBearerAuth(ctx, operationName, authcodegen.BearerAuth{
		Token: auth.Token,
		Roles: auth.Roles,
	})
	if err != nil {
		return nil, NewErrStorySchematicsSecurityHandler(err)
	}

	return ctx, nil
}
