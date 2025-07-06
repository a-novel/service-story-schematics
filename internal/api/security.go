package api

import (
	"context"
	"fmt"

	authModels "github.com/a-novel/service-authentication/models"
	authPkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/api/codegen"
)

type SecurityHandler struct {
	handler *authPkg.HandleBearerAuth[codegen.OperationName]
}

func NewSecurity(
	source authPkg.AuthenticateSource, permissions authModels.PermissionsConfig,
) (*SecurityHandler, error) {
	handler, err := authPkg.NewHandleBearerAuth[codegen.OperationName](source, permissions)
	if err != nil {
		return nil, fmt.Errorf("NewSecurity: %w", err)
	}

	return &SecurityHandler{handler: handler}, nil
}

func (security *SecurityHandler) HandleBearerAuth(
	ctx context.Context, operationName codegen.OperationName, auth codegen.BearerAuth,
) (context.Context, error) {
	handler, err := security.handler.HandleBearerAuth(ctx, operationName, &auth)
	if err != nil {
		return nil, fmt.Errorf("HandleBearerAuth: %w", err)
	}

	return handler, nil
}
