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

type ACL struct {
	contract *gwclient.Contract
}

func NewACL(client *network.FabricClient, contract string) (*ACL, error) {
	if client == nil || contract == "" {
		return nil, errors.New("invalid arguments")
	}

	acl := &ACL{
		contract: client.Channel("").GetContract(contract),
	}

	return acl, nil
}

func (acl *ACL) SetRoleAdmin(role []byte, adminRole []byte) error {
	_, err := acl.contract.SubmitTransaction("SetRoleAdmin")
	if err != nil {
		return utils.ParseTxError(err)
	}
	return nil
}

func (acl *ACL) GetRoleAdmin(role []byte) ([]byte, error) {
	result, err := acl.contract.EvaluateTransaction("GetRoleAdmin", string(role))
	if err != nil {
		return nil, utils.ParseTxError(err)
	}
	return result, nil
}

func (acl *ACL) HasRole(role []byte, account string) (string, error) {
	result, err := acl.contract.EvaluateTransaction("HasRole", string(role), account)
	if err != nil {
		return "", utils.ParseTxError(err)
	}
	return string(result), nil
}

func (acl *ACL) GrantRole(role []byte, account string) error {
	_, err := acl.contract.SubmitTransaction("GrantRole", string(role), account)
	if err != nil {
		return utils.ParseTxError(err)
	}
	return nil
}

func (acl *ACL) RevokeRole(role []byte, account string) error {
	_, err := acl.contract.SubmitTransaction("RevokeRole", string(role), account)
	if err != nil {
		return utils.ParseTxError(err)
	}
	return nil
}

func (acl *ACL) RenounceRole(msg *utils.Message, role []byte, account string) error {
	rawMsg, err := msg.Marshal()
	if err != nil {
		return err
	}
	_, err = acl.contract.SubmitTransaction("RenounceRole", string(rawMsg), string(role), account)
	if err != nil {
		return utils.ParseTxError(err)
	}
	return nil
}
