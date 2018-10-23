/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package login

import (
	"context"
	"github.com/nalej/grpc-authx-go"
)

// Manager structure with the required clients for roles operations.
type Manager struct {
	accessClient grpc_authx_go.AuthxClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(accessClient grpc_authx_go.AuthxClient) Manager {
	return Manager{accessClient:accessClient}
}

// LoginWithBasicCredentials performs the login of a user with a set of basic credentials. If the login
// is successful, it will return a JWT token.
func (m * Manager) LoginWithBasicCredentials(loginRequest *grpc_authx_go.LoginWithBasicCredentialsRequest) (*grpc_authx_go.LoginResponse, error){
	return m.accessClient.LoginWithBasicCredentials(context.Background(), loginRequest)
}