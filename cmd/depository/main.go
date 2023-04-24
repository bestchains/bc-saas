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
	"fmt"
	"os"

	"github.com/bestchains/bc-explorer/pkg/auth"
	"github.com/bestchains/bc-explorer/pkg/network"
	"github.com/bestchains/bc-saas/pkg/contracts"
	handler "github.com/bestchains/bc-saas/pkg/handlers"
	"github.com/bestchains/bc-saas/pkg/listener"
	"github.com/bestchains/bc-saas/pkg/models"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"k8s.io/klog/v2"
)

var (
	profile    = flag.String("profile", "./network.json", "profile to connect with blockchain network")
	contract   = flag.String("contract", "depository", "contract name")
	addr       = flag.String("addr", ":9999", "used to listen and serve http requests")
	db         = flag.String("db", "pg", "which database to use, default is pg(postgresql)")
	dsn        = flag.String("dsn", "postgres://bestchains:Passw0rd!@127.0.0.1:5432/bc-saas?sslmode=disable", "database connection string")
	authMethod = flag.String("auth", "none", "user authentication method, none, oidc or kubernetes")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		klog.Error(err)
	}
}

func run() error {
	// initialize contract client
	raw, err := os.ReadFile(*profile)
	if err != nil {
		return err
	}
	profile := &network.Network{}
	err = json.Unmarshal(raw, profile)
	if err != nil {
		return err
	}
	client, err := network.NewFabricClient(profile)
	if err != nil {
		return err
	}

	pctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher := listener.NewLogListener()
	klog.Info("init db...")

	klog.Infoln("Starting a digital depository server")

	klog.Infoln("Creating http server")
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		Immutable:     true,
		AppName:       "bc-saas",
	})

	app.Use(cors.New(cors.ConfigDefault))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Use(auth.New(context.TODO(), auth.Config{
		AuthMethod:    *authMethod,
		SkipAuthorize: true,
	}))
	depository := app.Group("depository")

	// hyperledger handlers
	hfContract, err := contracts.NewHyperledger(client, *contract)
	if err != nil {
		return err
	}
	hfHandler := handler.NewHyperledgerHandler(hfContract)
	// hyperledger routes
	hf := depository.Group("hf")
	hf.Get("metadata", hfHandler.GetMetadata)

	// basic handlers
	basicContract, err := contracts.NewBasic(client, *contract)
	if err != nil {
		return err
	}
	basicHandler := handler.NewBasicHandler(basicContract)
	// basic routes
	basic := depository.Group("basic")
	basic.Get("currentNonce", basicHandler.CurrentNonce)
	basic.Get("total", basicHandler.Total)
	basic.Post("putValue", basicHandler.PutValue)
	basic.Get("getValue", basicHandler.GetValue)
	basic.Post("verifyValue", basicHandler.VerifyValue)

	if *db == "pg" {
		klog.Infoln("Using postgreSQL")
		opts, err := pg.ParseURL(*dsn)
		if err != nil {
			return err
		}
		pgDB := pg.Connect(opts)
		defer pgDB.Close()
		if err := pgDB.Ping(pctx); err != nil {
			panic(err)
		}
		orm.SetTableNameInflector(func(s string) string {
			return fmt.Sprintf("%s_%s_%s", profile.ID, profile.Channel, s)
		})

		if err := models.Init(pgDB); err != nil {
			panic(err)
		}
		watcher, err = listener.NewListener(client, basicContract, pgDB, *contract, profile.Channel)
		if err != nil {
			panic(err)
		}
	}
	go watcher.Events(pctx)
	// NOTE: DISABLE ACL
	// acl handlers
	// aclContract, err := contracts.NewACL(client, *contract)
	// if err != nil {
	// 	return err
	// }
	// acl routes
	// aclHandler := handler.NewACLHandler(aclContract)
	// aclGroup := depository.Group("acl")
	// aclGroup.Get("hasRole", aclHandler.HasRole)

	if err := app.Listen(*addr); err != nil {
		return err
	}

	return nil
}
