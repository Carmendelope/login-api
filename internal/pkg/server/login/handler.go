/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package login

import (
	"context"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/login-api/internal/pkg/entities"
)

// Handler structure for the user requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

// LoginWithBasicCredentials performs the login of a user with a set of basic credentials. If the login
// is successful, it will return a JWT token.
func (h * Handler) LoginWithBasicCredentials(ctx context.Context, loginRequest *grpc_authx_go.LoginWithBasicCredentialsRequest) (*grpc_authx_go.LoginResponse, error){
	err := entities.ValidLoginWithBasicCredentialsRequest(loginRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.LoginWithBasicCredentials(loginRequest)
}