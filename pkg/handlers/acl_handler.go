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
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/sha3"
)

type Role string

const (
	RoleAdmin  Role = "role~admin"
	RoleClient Role = "role~client"
)

func (role Role) Hashed() []byte {
	digest := sha3.Sum256([]byte(role))
	return digest[:]
}

type Account struct {
	Role    Role   `json:"role,omitempty"`
	Address string `json:"address,omitempty"`
}

// TODO: add all access control handlers
type ACLHandler struct {
	acl *contracts.ACL
}

func NewACLHandler(acl *contracts.ACL) ACLHandler {
	return ACLHandler{
		acl: acl,
	}
}

func (handler *ACLHandler) HasRole(ctx *fiber.Ctx) error {
	account := Account{
		Role:    Role(ctx.Query("role")),
		Address: ctx.Query("account"),
	}
	if account.Role == "" || account.Address == "" {
		return fiber.NewError(fiber.StatusBadRequest, "empty role or account")
	}

	result, err := handler.acl.HasRole(account.Role.Hashed(), account.Address)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(result)
}
