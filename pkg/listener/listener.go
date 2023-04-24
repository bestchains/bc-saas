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

package listener

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bestchains/bc-explorer/pkg/network"
	"github.com/bestchains/bc-saas/pkg/contracts"
	handler "github.com/bestchains/bc-saas/pkg/handlers"
	"github.com/bestchains/bc-saas/pkg/models"
	"github.com/go-pg/pg/v10"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"k8s.io/klog/v2"
)

const (
	BasicChaincodeEvent          = "PutValue"
	BasicChaincodeUntrustedEvent = "PutUntrustValue"
	TestUntrustEventENV          = "PUTUNTRUSTVALUE"
)

type listener struct {
	contractName   string
	channel        string
	events         <-chan *client.ChaincodeEvent
	fabClient      *network.FabricClient
	contractClient *contracts.Basic
	db             *pg.DB
}

func NewListener(fabClient *network.FabricClient, cc *contracts.Basic, db *pg.DB, contractName, channel string) (Listener, error) {
	l := &listener{
		contractName:   contractName,
		channel:        channel,
		fabClient:      fabClient,
		contractClient: cc,
		db:             db,
	}
	ctx := context.Background()
	events, err := fabClient.Channel("").ChaincodeEvents(ctx, contractName)
	if err != nil {
		return nil, err
	}
	l.events = events
	return l, nil
}

type kv struct {
	Index    uint64 `json:"index,omitempty"`
	KID      string `json:"kid,omitempty"`
	Operator string `json:"operator"`
	Owner    string `json:"owner"`
}

func (l *listener) Events(ctx context.Context) {
	klog.Infof("starting fetch contract %s's events", l.contractName)
	whichMode := os.Getenv(TestUntrustEventENV)
	for {
		select {
		case e := <-l.events:
			if whichMode == "" && e.EventName != BasicChaincodeEvent {
				klog.Warningf("Event %s, expect %s, skip", e.EventName, BasicChaincodeEvent)
				continue
			}
			if whichMode != "" && e.EventName != BasicChaincodeUntrustedEvent {
				klog.Warningf("Event %s, expect %s, skip", e.EventName, BasicChaincodeUntrustedEvent)
				continue
			}

			klog.V(5).Infof("[Debug] event %+v", *e)

			eventPayload := kv{}
			if err := json.Unmarshal(e.Payload, &eventPayload); err != nil {
				klog.Errorf("[Error] unmarshal event payload error %s", err)
				continue
			}
			klog.V(5).Infof("[Debug] event payload %+v", eventPayload)

			what, err := l.contractClient.GetValueByKID(eventPayload.KID)
			if err != nil {
				klog.Errorf("[Error] GetValue By KID %s, error %s", eventPayload.KID, err)
				continue
			}
			klog.V(5).Infof("[Debug] Call GetValueByKID %s get %s", eventPayload.KID, what)

			vdBytes, err := base64.StdEncoding.DecodeString(what)
			if err != nil {
				klog.Errorf("[Error] decode value %s erorr %s", what, err)
				continue
			}
			klog.V(5).Infof("[Debug] value depository bytes: %s", string(vdBytes))

			vd := handler.ValueDepository{}
			if err := json.Unmarshal(vdBytes, &vd); err != nil {
				klog.Errorf("[Error] unmarshal valueDepository error %s", err)
				continue
			}

			d := models.Depository{
				Index:            fmt.Sprintf("%d", eventPayload.Index),
				KID:              eventPayload.KID,
				Platform:         vd.Platform,
				Operator:         eventPayload.Operator,
				Owner:            eventPayload.Owner,
				BlockNumber:      fmt.Sprintf("%d", e.BlockNumber),
				ContentName:      vd.Name,
				ContentID:        vd.ContentID,
				ContentType:      vd.ContentID,
				TrustedTimestamp: vd.TrustedTimestamp,
			}
			klog.V(5).Infof("[Debug] insert vd %+v, d: %+v into db", vd, d)

			if _, err := l.db.Model(&d).Insert(); err != nil {
				klog.Errorf("[Error] failed to insert data, return error %s", err)
				continue
			}
			klog.Infof("[Success] insert %s to db", d.BlockNumber)
		case <-ctx.Done():
			klog.Info("context break down")
			return
		}
	}
}
