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
	"github.com/bestchains/bc-explorer/pkg/network"
	"github.com/bestchains/bc-saas/pkg/utils"
	gwclient "github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/pkg/errors"
)

type Hyperledger struct {
	contract *gwclient.Contract
}

func NewHyperledger(client *network.FabricClient, contract string) (*Hyperledger, error) {
	if client == nil {
		return nil, errors.New("invalid arguments")
	}

	acl := &Hyperledger{
		contract: client.Channel("").GetContract(contract),
	}

	return acl, nil
}

func (hf *Hyperledger) GetMetadata() ([]byte, error) {
	result, err := hf.contract.EvaluateTransaction("org.hyperledger.fabric:GetMetadata")
	if err != nil {
		return nil, utils.ParseTxError(err)
	}
	return result, nil
}
