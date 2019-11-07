/*
 * Copyright 2019 Nalej
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package login

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/login-api/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

// Handler structure for the user requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

// LoginWithBasicCredentials performs the login of a user with a set of basic credentials. If the login
// is successful, it will return a JWT token.
func (h *Handler) LoginWithBasicCredentials(ctx context.Context, loginRequest *grpc_authx_go.LoginWithBasicCredentialsRequest) (*grpc_authx_go.LoginResponse, error) {
	err := entities.ValidLoginWithBasicCredentialsRequest(loginRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	response, lgErr := h.Manager.LoginWithBasicCredentials(loginRequest)
	if lgErr != nil {
		log.Error().Str("trace", conversions.ToDerror(lgErr).DebugReport()).Msg("login error")
		return nil, conversions.ToGRPCError(derrors.NewGenericError("Invalid credentials"))
	}

	return response, lgErr
}
