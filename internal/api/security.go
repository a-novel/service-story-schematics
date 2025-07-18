package api

import (
	"context"
	"fmt"

	authconfig "github.com/a-novel/service-authentication/models/config"
	authpkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/models/api"
)

type SecurityHandler struct {
	handler *authpkg.HandleBearerAuth[apimodels.OperationName]
}

func NewSecurity(
	source authpkg.AuthenticateSource, permissions authconfig.Permissions,
) (*SecurityHandler, error) {
	handler, err := authpkg.NewHandleBearerAuth[apimodels.OperationName](source, permissions)
	if err != nil {
		return nil, fmt.Errorf("NewSecurity: %w", err)
	}

	return &SecurityHandler{handler: handler}, nil
}

func (security *SecurityHandler) HandleBearerAuth(
	ctx context.Context, operationName apimodels.OperationName, auth apimodels.BearerAuth,
) (context.Context, error) {
	handler, err := security.handler.HandleBearerAuth(ctx, operationName, &auth)
	if err != nil {
		return nil, fmt.Errorf("HandleBearerAuth: %w", err)
	}

	return handler, nil
}
