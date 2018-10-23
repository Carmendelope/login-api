/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-login-api-go"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-utils/pkg/tools"
	"github.com/nalej/login-api/internal/pkg/server/login"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
)

type Service struct {
	Configuration Config
	Server * tools.GenericGRPCServer
}

// NewService creates a new system model service.
func NewService(conf Config) *Service {
	return &Service{
		conf,
		tools.NewGenericGRPCServer(uint32(conf.Port)),
	}
}

type Clients struct {
	accessClient grpc_authx_go.AuthxClient
}

func (s * Service) GetClients() (* Clients, derrors.Error) {
	authxConn, err := grpc.Dial(s.Configuration.AuthxAddress, grpc.WithInsecure())
	if err != nil{
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

func (s * Service) LaunchHTTP() error {
	addr := fmt.Sprintf(":%d", s.Configuration.HTTPPort)
	clientAddr := fmt.Sprintf(":%d", s.Configuration.Port)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	mux := runtime.NewServeMux()

	if err := grpc_login_api_go.RegisterLoginHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start organizations handler")
	}

	log.Info().Str("address", addr).Msg("HTTP Listening")
	return http.ListenAndServe(addr, mux)
}

func (s * Service) LaunchGRPC() error {
	clients, cErr := s.GetClients()
	if cErr != nil{
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