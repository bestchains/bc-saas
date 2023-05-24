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

package handler

import (
	"github.com/bestchains/bc-saas/pkg/contracts"
	"github.com/bestchains/bc-saas/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

// KeyValue defines common key-value fields for a depository
type Repository struct {
	ID string `json:"id"`
	// URL of this repository
	URL string `json:"url,omitempty"`
	// Message is a base64 encoded string of utils.Message
	Message string `json:"message,omitempty"`
}

type MarketHandler struct {
	market *contracts.Market
}

// NewMarketHandler creates a new instance of MarketHandler.
//
// It takes a pointer to a Market contract and returns a pointer to a MarketHandler.
func NewMarketHandler(contractClient *contracts.Market) *MarketHandler {
	return &MarketHandler{
		market: contractClient,
	}
}

// CurrentNonce returns the current nonce for the given account.
func (lh *MarketHandler) CurrentNonce(ctx *fiber.Ctx) error {
	account := ctx.Query("account")
	nonce, err := lh.market.CurrentNonce(account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(fiber.Map{
		"nonce": nonce,
	})
}

// CreateRepo creates a new repository with a market in the given URL.
//
// ctx: Context from the request.
// Returns an error if there is a bad request or internal server error.
func (lh *MarketHandler) CreateRepo(ctx *fiber.Ctx) error {
	var err error

	// Parse request body into a new Repository struct
	repo := new(Repository)
	err = ctx.BodyParser(repo)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Check that URL and message are not empty
	if repo.URL == "" || repo.Message == "" {
		return fiber.NewError(fiber.StatusBadRequest, "url and message cannot be empty")
	}

	// Unmarshal the message from base64 string
	message := new(utils.Message)
	if err = message.UnmarshalBase64Str(repo.Message); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, errors.Wrap(err, "invalid message").Error())
	}

	// Create the repository in the market
	repoID, err := lh.market.CreateRepo(message, repo.URL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return the repository ID as JSON
	return ctx.JSON(fiber.Map{
		"repo_id": repoID,
	})
}

// UpdateRepo updates the repository with the given ID and URL.
// It expects a JSON-encoded body containing a Repository object.
// Returns a JSON response with the updated repository ID and URL.
func (lh *MarketHandler) UpdateRepo(ctx *fiber.Ctx) error {
	var err error
	// Parse request body into Repository object
	repo := new(Repository)
	err = ctx.BodyParser(repo)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Check that required fields are not empty
	if repo.ID == "" || repo.URL == "" || repo.Message == "" {
		return fiber.NewError(fiber.StatusBadRequest, "id,url&message cannot be empty")
	}

	// Unmarshal Base64-encoded message string into a Message object
	message := new(utils.Message)
	if err = message.UnmarshalBase64Str(repo.Message); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, errors.Wrap(err, "invalid message").Error())
	}

	// Update repository in market
	err = lh.market.UpdateRepo(message, repo.ID, repo.URL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return updated repository ID and URL as JSON
	return ctx.JSON(fiber.Map{
		"repo_id": repo.ID,
		"url":     repo.URL,
	})
}

// GetRepos returns a list of repositories from the MarketHandler's market instance
func (lh *MarketHandler) GetRepos(ctx *fiber.Ctx) error {
	// Get the list of repositories from the market instance
	repos, err := lh.market.GetRepos()
	if err != nil {
		// Return an error response if there was an error getting the repositories
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return the list of repositories as a JSON response
	return ctx.JSON(repos)
}
