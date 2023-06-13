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

package events

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"

	"github.com/bestchains/bc-saas/pkg/contracts"
	handler "github.com/bestchains/bc-saas/pkg/handlers"
	"github.com/bestchains/bc-saas/pkg/models"
)

const (
	DepositoryEventPutValue        Event = "PutValue"
	DepositoryEventPutUntrustValue Event = "PutUntrustValue"
)

type kv struct {
	Index    uint64 `json:"index,omitempty"`
	KID      string `json:"kid,omitempty"`
	Operator string `json:"operator"`
	Owner    string `json:"owner"`
}

type DepositoryEventHandler struct {
	contractClient *contracts.Depository
	db             *pg.DB
}

func NewDepositoryEventHandler(contractClient *contracts.Depository, db *pg.DB) *DepositoryEventHandler {
	return &DepositoryEventHandler{
		contractClient: contractClient,
		db:             db,
	}
}

// HandlePutValue handle events
// - EventPutValue
// - EventPutUntrustValue
func (deh *DepositoryEventHandler) HandlePutValue(e *client.ChaincodeEvent) error {
	eventPayload := kv{}
	if err := json.Unmarshal(e.Payload, &eventPayload); err != nil {
		return errors.Wrap(err, "unmarshal event payload")
	}
	klog.V(5).Infof("[Debug] event payload %+v", eventPayload)

	what, err := deh.contractClient.GetValueByKID(eventPayload.KID)
	if err != nil {
		return errors.Wrapf(err, "getValue By KID %s", eventPayload.KID)
	}

	klog.V(5).Infof("[Debug] Call GetValueByKID %s get %s", eventPayload.KID, what)

	vdBytes, err := base64.StdEncoding.DecodeString(what)
	if err != nil {
		return errors.Wrap(err, "decode value")
	}
	klog.V(5).Infof("[Debug] value depository bytes: %s", string(vdBytes))

	vd := handler.ValueDepository{}
	if err := json.Unmarshal(vdBytes, &vd); err != nil {
		return errors.Wrap(err, "unmarshal valueDepository")
	}

	d := models.Depository{
		Index:            fmt.Sprintf("%d", eventPayload.Index),
		KID:              eventPayload.KID,
		Platform:         vd.Platform,
		Operator:         eventPayload.Operator,
		Owner:            eventPayload.Owner,
		BlockNumber:      e.BlockNumber,
		Name:             vd.Name,
		ContentName:      vd.ContentName,
		ContentID:        vd.ContentID,
		ContentType:      vd.ContentID,
		TrustedTimestamp: fmt.Sprintf("%d", time.Now().Unix()),
		Description:      vd.Description,
	}
	klog.V(5).Infof("[Debug] insert vd %+v, d: %+v into db", vd, d)

	if _, err := deh.db.Model(&d).Insert(); err != nil {
		return errors.Wrap(err, "insert depository data")
	}
	klog.Infof("[Success] insert depository %s at block %d to db", d.KID, d.BlockNumber)
	return nil
}
