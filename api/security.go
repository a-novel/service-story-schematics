package api

import (
	"errors"

	authapi "github.com/a-novel/authentication/api"
	authcodegen "github.com/a-novel/authentication/api/codegen"
	"github.com/a-novel/authentication/models"

	"github.com/a-novel-kit/context"

	"github.com/a-novel/story-schematics/api/codegen"
)

var ErrStorySchematicsSecurityHandler = errors.New("story schematics security handler")

func NewErrStorySchematicsSecurityHandler(err error) error {
	return errors.Join(err, ErrStorySchematicsSecurityHandler)
}

type SecurityHandler struct {
	authHandler *authapi.SecurityHandler
}

func (security *SecurityHandler) HandleBearerAuth(
	ctx context.Context, operationName codegen.OperationName, auth codegen.BearerAuth,
) (context.Context, error) {
	ctx, err := security.authHandler.HandleBearerAuth(ctx, operationName, authcodegen.BearerAuth{
		Token: auth.Token,
	})
	if err != nil {
		return nil, NewErrStorySchematicsSecurityHandler(err)
	}

	return ctx, nil
}

func NewSecurity(
	required map[codegen.OperationName][]models.Permission,
	granted models.PermissionsConfig,
	authService authapi.SecurityHandlerService,
) (*SecurityHandler, error) {
	authHandler, err := authapi.NewSecurity(required, granted, authService)
	if err != nil {
		return nil, NewErrStorySchematicsSecurityHandler(err)
	}

	return &SecurityHandler{
		authHandler: authHandler,
	}, nil
}
