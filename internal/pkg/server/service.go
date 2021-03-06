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

package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-login-api-go"
	"github.com/nalej/login-api/internal/pkg/server/login"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"strings"
)

type Service struct {
	Configuration Config
}

// NewService creates a new system model service.
func NewService(conf Config) *Service {
	return &Service{
		conf,
	}
}

type Clients struct {
	accessClient grpc_authx_go.AuthxClient
}

func (s *Service) GetClients() (*Clients, derrors.Error) {
	authxConn, err := grpc.Dial(s.Configuration.AuthxAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the authx manager")
	}

	aClient := grpc_authx_go.NewAuthxClient(authxConn)

	return &Clients{aClient}, nil
}

// Run the service, launch the REST service handler.
func (s *Service) Run() error {
	s.Configuration.Print()

	go s.LaunchGRPC()
	return s.LaunchHTTP()
}

// allowCORS allows Cross Origin Resource Sharing from any origin.
// Don't do this without consideration in production systems.
func (s *Service) allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	log.Debug().Str("URL", r.URL.Path).Msg("preflight request")
}

func (s *Service) LaunchHTTP() error {
	addr := fmt.Sprintf(":%d", s.Configuration.HTTPPort)
	clientAddr := fmt.Sprintf(":%d", s.Configuration.Port)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	mux := runtime.NewServeMux()

	if err := grpc_login_api_go.RegisterLoginHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start organizations handler")
	}

	server := &http.Server{
		Addr:    addr,
		Handler: s.allowCORS(mux),
	}

	log.Info().Str("address", addr).Msg("HTTP Listening")
	return server.ListenAndServe()
	//return http.ListenAndServe(addr, mux)
}

func (s *Service) LaunchGRPC() error {
	clients, cErr := s.GetClients()
	if cErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("cannot generate clients")
		return cErr
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
	}

	// Create handlers
	loginManager := login.NewManager(clients.accessClient)
	loginHandler := login.NewHandler(loginManager)

	grpcServer := grpc.NewServer()
	grpc_login_api_go.RegisterLoginServer(grpcServer, loginHandler)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}
