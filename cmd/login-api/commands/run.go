/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// This is an example of an executable command.

package commands

import (
	"github.com/nalej/login-api/internal/pkg/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var config = server.Config{}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launch the server API",
	Long:  `Launch the server API`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		log.Info().Msg("Launching API!")
		server := server.NewService(config)
		server.Run()
	},
}

func init() {
	runCmd.Flags().IntVar(&config.Port, "port", 8444, "Port to launch the Public gRPC API")
	runCmd.Flags().IntVar(&config.HTTPPort, "httpPort", 8443, "Port to launch the Public HTTP API")
	runCmd.PersistentFlags().StringVar(&config.AuthxAddress, "authxAddress", "localhost:8810",
		"Authx address (host:port)")
	rootCmd.AddCommand(runCmd)
}