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
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-login-api-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"math/rand"
	"os"
)

var _ = ginkgo.Describe("Applications", func() {

	var runIntegration = os.Getenv("RUN_INTEGRATION_TEST")

	if runIntegration != "true" {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		authxAddress = os.Getenv("IT_AUTHX_ADDRESS")
	)

	if authxAddress == "" {
		ginkgo.Fail("missing environment variables")
	}

	var authxClient grpc_authx_go.AuthxClient
	var server *grpc.Server
	var listener *bufconn.Listener
	var client grpc_login_api_go.LoginClient

	var role grpc_authx_go.Role
	var credentials grpc_authx_go.AddBasicCredentialRequest

	ginkgo.BeforeSuite(func() {

		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		authxConn, err := grpc.Dial(authxAddress, grpc.WithInsecure())
		gomega.Expect(err).To(gomega.Succeed())

		authxClient = grpc_authx_go.NewAuthxClient(authxConn)

		manager := NewManager(authxClient)
		handler := NewHandler(manager)

		grpc_login_api_go.RegisterLoginServer(server, handler)
		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_login_api_go.NewLoginClient(conn)
		rand.Seed(ginkgo.GinkgoRandomSeed())

		// create a user (credentials and role)
		role = grpc_authx_go.Role{
			OrganizationId: "orgID",
			RoleId:         "roleId",
			Name:           "rName1",
			Primitives:     []grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG},
			Internal:       false,
		}

		authxClient.AddRole(context.Background(), &role)

		credentials = grpc_authx_go.AddBasicCredentialRequest{
			OrganizationId: role.OrganizationId,
			Username:       "email@nalej.com",
			Password:       "password",
			RoleId:         role.RoleId,
		}

		authxClient.AddBasicCredentials(context.Background(), &credentials)

	})

	ginkgo.AfterSuite(func() {

		// TODO: remove credentials and role (when the methods are added)

		server.Stop()
		listener.Close()

	})

	ginkgo.Context("login", func() {
		ginkgo.It("should not be able to login, user not found", func() {

			toLogin := grpc_authx_go.LoginWithBasicCredentialsRequest{
				Username: "wrongEmail@nalej.com",
				Password: credentials.Password,
			}
			_, err := client.LoginWithBasicCredentials(context.Background(), &toLogin)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("should not be able to login, error in password", func() {

			toLogin := grpc_authx_go.LoginWithBasicCredentialsRequest{
				Username: credentials.Username,
				Password: "wrongPassword",
			}
			_, err := client.LoginWithBasicCredentials(context.Background(), &toLogin)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("should be able to login", func() {

			toLogin := grpc_authx_go.LoginWithBasicCredentialsRequest{
				Username: credentials.Username,
				Password: credentials.Password,
			}
			response, err := client.LoginWithBasicCredentials(context.Background(), &toLogin)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(response).NotTo(gomega.BeNil())
		})
	})
})
