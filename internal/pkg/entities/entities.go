/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
)

func ValidLoginWithBasicCredentialsRequest(loginRequest *grpc_authx_go.LoginWithBasicCredentialsRequest) derrors.Error{
	if loginRequest.Username == "" {
		return derrors.NewInvalidArgumentError("expecting username")
	}
	if loginRequest.Password == "" {
		return derrors.NewInvalidArgumentError("expecting password")
	}
	return nil
}