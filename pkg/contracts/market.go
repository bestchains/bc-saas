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

type Market struct {
	contract *gwclient.Contract
}

// NewMarket creates a new instance of Market using the provided FabricClient and contract name.
// Returns a pointer to Market and an error if the arguments are invalid or the contract is not found.
func NewMarket(client *network.FabricClient, contract string) (*Market, error) {
	// Check that the arguments are valid
	if client == nil || contract == "" {
		return nil, errors.New("invalid arguments")
	}

	// Create a new instance of Market
	market := &Market{
		contract: client.Channel("").GetContract(contract),
	}

	return market, nil
}

// Initialize initializes the market.
func (market *Market) Initialize() error {
	// TODO: Implement initialization logic.
	return nil
}

// CurrentNonce returns the current nonce for a given account on the market contract.
func (market *Market) CurrentNonce(account string) (uint64, error) {
	// Evaluate the "Current" transaction on the contract with the given account.
	result, err := market.contract.EvaluateTransaction("Current", account)
	if err != nil {
		// If an error occurred, parse it and return it.
		return 0, utils.ParseTxError(err)
	}

	// Otherwise, parse the result as a uint64 and return it.
	return strconv.ParseUint(string(result), 10, 64)
}

// CreateRepo creates a new repository with the specified message and URL.
// It submits a new transaction to the blockchain via the smart contract.
// Returns the ID of the new repository or an error if the transaction fails.
func (market *Market) CreateRepo(msg *utils.Message, url string) (string, error) {
	// Marshal the message to raw bytes
	rawMsg, err := msg.Marshal()
	if err != nil {
		return "", err
	}

	// Submit the transaction to the smart contract
	repoID, err := market.contract.SubmitTransaction("CreateRepo", string(rawMsg), url)
	if err != nil {
		return "", utils.ParseTxError(err)
	}

	return string(repoID), nil
}

// UpdateRepo updates the URL for a given repository ID in the blockchain
func (market *Market) UpdateRepo(msg *utils.Message, repoID string, newUrl string) error {
	// Marshal the message to bytes
	rawMsg, err := msg.Marshal()
	if err != nil {
		return err
	}

	// Call the "UpdateRepo" function on the blockchain contract
	_, err = market.contract.SubmitTransaction("UpdateRepo", string(rawMsg), repoID, newUrl)
	if err != nil {
		return utils.ParseTxError(err)
	}

	// Return nil if everything worked as expected
	return nil
}

// GetRepos returns the repositories associated with the market.
func (market *Market) GetRepos() ([]byte, error) {
	// Evaluate the "GetRepos" transaction on the contract.
	result, err := market.contract.EvaluateTransaction("GetRepos")
	if err != nil {
		// If there was an error, return it as a parsed transaction error.
		return nil, utils.ParseTxError(err)
	}
	// Otherwise, return the result.
	return result, nil
}
