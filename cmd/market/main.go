/*
Copyright 2023 The Bestchains Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"

	"github.com/bestchains/bc-explorer/pkg/auth"
	"github.com/bestchains/bc-explorer/pkg/network"
	"github.com/bestchains/bc-saas/pkg/contracts"
	handler "github.com/bestchains/bc-saas/pkg/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"k8s.io/klog/v2"
)

var (
	profile     = flag.String("profile", "./network.json", "profile to connect with blockchain network")
	contract    = flag.String("contract", "market", "contract name")
	addr        = flag.String("addr", ":9998", "used to listen and serve http requests")
	authMethod  = flag.String("auth", "none", "user authentication method, none, oidc or kubernetes")
	enablePprof = flag.Bool("enable-pprof", false, "enable performance profiling in depository service")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		klog.Error(err)
	}
}

// run starts the server and initializes the contract client
func run() error {
	// initialize contract client
	raw, err := os.ReadFile(*profile) // read the profile file
	if err != nil {
		return err
	}
	profile := &network.Network{}
	err = json.Unmarshal(raw, profile) // unmarshal the profile data into a Network struct
	if err != nil {
		return err
	}
	fabClient, err := network.NewFabricClient(profile) // create a new Fabric client
	if err != nil {
		return err
	}

	klog.Info("init contract client")
	contractClient, err := contracts.NewMarket(fabClient, *contract) // create a new Market contract client
	if err != nil {
		return err
	}

	klog.Infoln("Creating http server")
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		Immutable:     true,
		AppName:       "bc-saas",
	})

	app.Use(cors.New(cors.ConfigDefault)) // add CORS middleware
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	})) // add logger middleware
	app.Use(auth.New(context.TODO(), auth.Config{
		AuthMethod:    *authMethod,
		SkipAuthorize: true,
	})) // add auth middleware

	//  Enable pprof
	if *enablePprof {
		app.Use(pprof.New())
	}

	// hyperledger handlers
	hfContract, err := contracts.NewHyperledger(fabClient, *contract) // create a new Hyperledger contract client
	if err != nil {
		return err
	}
	hfHandler := handler.NewHyperledgerHandler(hfContract) // create a new Hyperledger handler
	// hyperledger routes
	hf := app.Group("hf")                     // create a new route group for Hyperledger
	hf.Get("metadata", hfHandler.GetMetadata) // add a route to get metadata

	licenseHandler := handler.NewMarketHandler(contractClient) // create a new Market handler

	// market routes
	license := app.Group("market")                    // create a new route group for Market
	license.Get("nonce", licenseHandler.CurrentNonce) // add a route to get the current nonce
	license.Post("repo", licenseHandler.CreateRepo)   // add a route to create a new repo
	license.Put("repo", licenseHandler.UpdateRepo)    // add a route to update a repo
	license.Get("repos", licenseHandler.GetRepos)     // add a route to get all repos

	klog.Infoln("Starting a market server")

	if err := app.Listen(*addr); err != nil { // start the server
		return err
	}

	return nil
}
