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

package contracts

import (
	"strconv"

	"github.com/bestchains/bc-explorer/pkg/network"
	"github.com/bestchains/bc-saas/pkg/utils"
	gwclient "github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/pkg/errors"
)

// Code Generator
type Basic struct {
	contract *gwclient.Contract
}

func NewBasic(client *network.FabricClient, contract string) (*Basic, error) {
	if client == nil || contract == "" {
		return nil, errors.New("invalid arguments")
	}

	basic := &Basic{
		contract: client.Channel("").GetContract(contract),
	}

	return basic, nil
}

func (basic *Basic) Initialize() error {
	_, err := basic.contract.SubmitTransaction("Initialize")
	if err != nil {
		return utils.ParseTxError(err)
	}
	return nil
}

func (basic *Basic) Total() (uint64, error) {
	result, err := basic.contract.EvaluateTransaction("Total")
	if err != nil {
		return 0, utils.ParseTxError(err)
	}
	return strconv.ParseUint(string(result), 10, 64)
}

func (basic *Basic) PutValue(msg *utils.Message, val string) (string, error) {
	rawMsg, err := msg.Marshal()
	if err != nil {
		return "", err
	}
	kid, err := basic.contract.SubmitTransaction("PutValue", string(rawMsg), val)
	if err != nil {
		return "", utils.ParseTxError(err)
	}

	return string(kid), nil
}

func (basic *Basic) GetValueByIndex(index string) (string, error) {
	result, err := basic.contract.EvaluateTransaction("GetValueByIndex", index)
	if err != nil {
		return "", utils.ParseTxError(err)
	}
	return string(result), nil
}

func (basic *Basic) GetValueByKID(kid string) (string, error) {
	result, err := basic.contract.EvaluateTransaction("GetValueByKID", kid)
	if err != nil {
		return "", utils.ParseTxError(err)
	}
	return string(result), nil
}
