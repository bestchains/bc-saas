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
	"github.com/bestchains/bc-saas/pkg/depositories"
	"github.com/bestchains/bc-saas/pkg/events"
	handler "github.com/bestchains/bc-saas/pkg/handlers"
	"github.com/bestchains/bc-saas/pkg/listener"
	"github.com/bestchains/bc-saas/pkg/models"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"k8s.io/klog/v2"
)

var (
	profile     = flag.String("profile", "./network.json", "profile to connect with blockchain network")
	contract    = flag.String("contract", "depository", "contract name")
	addr        = flag.String("addr", ":9999", "used to listen and serve http requests")
	db          = flag.String("db", "pg", "which database to use, default is pg(postgresql)")
	dsn         = flag.String("dsn", "postgres://bestchains:Passw0rd!@127.0.0.1:5432/bc-saas?sslmode=disable", "database connection string")
	authMethod  = flag.String("auth", "none", "user authentication method, none, oidc or kubernetes")
	enablePprof = flag.Bool("enable-pprof", false, "enable performance profiling in depository service")

	// flags for depository certificate generation
	templateImagePath = flag.String("cert-template-image", "resource/certificate_template.jpg", "template image for depository's certificate generation")
	ttfFontPath       = flag.String("cert-ttf-font", "resource/ttf/SourceHanSansCN-Normal.ttf", "ttf font file for depository's certificate generation")
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
	fabClient, err := network.NewFabricClient(profile)
	if err != nil {
		return err
	}

	pctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher := listener.NewLogListener()
	dbHandler := depositories.NewLoggerHandler()

	// basic handlers
	klog.Info("init contract client")
	contractClient, err := contracts.NewDepository(fabClient, *contract)
	if err != nil {
		return err
	}

	klog.Info("init db...")
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
		pgDB.AddQueryHook(&models.Depository{})
		orm.SetTableNameInflector(func(s string) string {
			return fmt.Sprintf("%s_%s_%s", profile.ID, profile.Channel, s)
		})

		if err := models.Init(pgDB); err != nil {
			panic(err)
		}

		dbHandler, err = depositories.NewDBHandler(pgDB, *templateImagePath, *ttfFontPath)
		if err != nil {
			panic(err)
		}
		// inject events to database once pg is used
		eventSub, err := fabClient.Channel(profile.Channel).ChaincodeEvents(pctx, *contract, client.WithStartBlock(
			models.MaxBlockNumber(pgDB),
		))
		if err != nil {
			panic(err)
		}
		// register Depository related events
		eventHandler := events.NewDepositoryEventHandler(contractClient, pgDB)
		watcher, err = listener.NewListener(eventSub, map[events.Event]events.EventHandler{
			events.DepositoryEventPutUntrustValue: eventHandler.HandlePutValue,
			events.DepositoryEventPutValue:        eventHandler.HandlePutValue,
		})
		if err != nil {
			panic(err)
		}
	}

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

	//  Enable pprof
	if *enablePprof {
		app.Use(pprof.New())
	}

	// hyperledger handlers
	hfContract, err := contracts.NewHyperledger(fabClient, *contract)
	if err != nil {
		return err
	}
	hfHandler := handler.NewHyperledgerHandler(hfContract)
	// hyperledger routes
	hf := app.Group("hf")
	hf.Get("metadata", hfHandler.GetMetadata)

	basicHandler := handler.NewBasicHandler(contractClient, dbHandler)
	// basic routes
	basic := app.Group("basic")
	basic.Get("currentNonce", basicHandler.CurrentNonce)
	basic.Get("total", basicHandler.Total)
	basic.Post("putValue", basicHandler.PutValue)
	basic.Post("putUntrustValue", basicHandler.PutUntrustValue)
	basic.Get("getValue", basicHandler.GetValue)
	basic.Post("verifyValue", basicHandler.VerifyValue)
	basic.Get("depositories", basicHandler.List)
	basic.Get("depositories/:kid", basicHandler.Get)
	basic.Get("depositories/certificate/:kid", basicHandler.GetDepositoryCertificate)

	klog.Infoln("Starting a digital depository server")

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
