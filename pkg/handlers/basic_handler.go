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
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/bestchains/bc-saas/pkg/contracts"
	"github.com/bestchains/bc-saas/pkg/depositories"
	"github.com/bestchains/bc-saas/pkg/utils"
	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

// KeyValue defines common key-value fields for a depository
type KeyValue struct {
	Index string `json:"index,omitempty"`
	KID   string `json:"kid,omitempty"`
	Value string `json:"value,omitempty"`
	// Message is a base64 encoded string of utils.Message
	Message string `json:"message,omitempty"`
}

// ValueDepository defines valuable fields for a depository
type ValueDepository struct {
	Name        string `json:"name"`
	ContentName string `json:"contentName"`
	ContentType string `json:"contentType"`
	ContentID   string `json:"contentID"` // hash of the file

	// ContentSize the size of file. unit is byte
	ContentSize      int64  `json:"contentSize"`
	TrustedTimestamp string `json:"trustedTimestamp"`
	Platform         string `json:"platform"`
	Description      string `json:"description,omitempty"`
}

// VerifyStatus defines response fields for a depository verification
type VerifyStatus struct {
	Status bool   `json:"status"`
	Reason string `json:"reason"`
}

type BasicHandler struct {
	contractClient *contracts.Depository
	dbHandler      depositories.Interface
}

func NewBasicHandler(contractClient *contracts.Depository, h depositories.Interface) BasicHandler {
	return BasicHandler{
		contractClient: contractClient,
		dbHandler:      h,
	}
}

func (h *BasicHandler) CurrentNonce(ctx *fiber.Ctx) error {
	account := ctx.Query("account")
	nonce, err := h.contractClient.CurrentNonce(account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(map[string]interface{}{
		"nonce": nonce,
	})
}

// Total
func (h *BasicHandler) Total(ctx *fiber.Ctx) error {
	total, err := h.contractClient.Total()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(map[string]interface{}{
		"total": total,
	})
}

func (h *BasicHandler) PutUntrustValue(ctx *fiber.Ctx) error {
	var err error
	// only value is required in untrust put
	kv := new(KeyValue)
	err = ctx.BodyParser(kv)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// validate if value is a `ValueDepository`
	if kv.Value == "" {
		return fiber.NewError(fiber.StatusBadRequest, "value cannot be empty")
	}
	rawValue, err := base64.StdEncoding.DecodeString(kv.Value)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, errors.Wrap(err, "invalid value").Error())
	}
	value := new(ValueDepository)
	if err = json.Unmarshal(rawValue, value); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, errors.Wrap(err, "invalid value").Error())
	}

	kid, err := h.contractClient.PutUntrustValue(kv.Value)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(&KeyValue{
		KID: kid,
	})
}

func (h *BasicHandler) PutValue(ctx *fiber.Ctx) error {
	var err error

	kv := new(KeyValue)
	err = ctx.BodyParser(kv)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// validate if value is a `ValueDepository`
	if kv.Value == "" {
		return fiber.NewError(fiber.StatusBadRequest, "value cannot be empty")
	}
	rawValue, err := base64.StdEncoding.DecodeString(kv.Value)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, errors.Wrap(err, "invalid value").Error())
	}
	value := new(ValueDepository)
	if err = json.Unmarshal(rawValue, value); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, errors.Wrap(err, "invalid value").Error())
	}

	// validate message
	if kv.Message == "" {
		return fiber.NewError(fiber.StatusBadRequest, "message cannot be empty")
	}
	message := new(utils.Message)
	if err = message.UnmarshalBase64Str(kv.Message); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, errors.Wrap(err, "invalid message").Error())
	}

	kid, err := h.contractClient.PutValue(message, kv.Value)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(&KeyValue{
		KID: kid,
	})

}

func (h *BasicHandler) VerifyValue(ctx *fiber.Ctx) error {
	var err error

	kvArgs := new(KeyValue)
	err = ctx.BodyParser(kvArgs)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// get kv with index or kid
	kv, err := h.getValue(*kvArgs)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// validate value
	if kv.Value != kvArgs.Value {
		return ctx.JSON(&VerifyStatus{
			Status: false,
			Reason: "value mismatch",
		})
	}

	return ctx.JSON(&VerifyStatus{
		Status: true,
		Reason: "value match",
	})
}

func (h *BasicHandler) GetValue(ctx *fiber.Ctx) error {
	arg := KeyValue{
		Index: ctx.Query("index"),
		KID:   ctx.Query("kid"),
	}

	if arg.Index == "" && arg.KID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "must provide depository index or kid")
	}

	kv, err := h.getValue(arg)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(kv)
}

func (h *BasicHandler) getValue(kv KeyValue) (*KeyValue, error) {
	if kv.Index != "" {
		value, err := h.contractClient.GetValueByIndex(kv.Index)
		if err != nil {
			return nil, err
		}
		kv.Value = value
	} else if kv.KID != "" {
		value, err := h.contractClient.GetValueByKID(kv.KID)
		if err != nil {
			return nil, err
		}
		kv.Value = value
	}

	return &kv, nil
}

func (h *BasicHandler) List(ctx *fiber.Ctx) error {
	klog.Info("BasicHandler List Depositories")
	klog.V(5).Infof(" with ctx %+v\n", *ctx)

	arg := depositories.DepositoryCond{
		From:        ctx.QueryInt("from", 0),
		Size:        ctx.QueryInt("size", 10),
		StartTime:   int64(ctx.QueryInt("startTime", 0)),
		EndTime:     int64(ctx.QueryInt("endTime", 0)),
		Name:        ctx.Query("name"),
		KID:         ctx.Query("kid"),
		ContentName: ctx.Query("contentName", ""),
	}

	result, count, err := h.dbHandler.List(arg)
	if err != nil {
		klog.Errorf("[Error] list depositories error %s", err)
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(map[string]string{
			"msg": err.Error(),
		})
	}
	data := map[string]interface{}{
		"data":  result,
		"count": count,
	}
	return ctx.JSON(data)
}

func (h *BasicHandler) Get(ctx *fiber.Ctx) error {
	klog.Info("BasicHandler Get Depository")
	klog.V(5).Infof(" with ctx %+v\n", *ctx)

	kid := ctx.Params("kid")
	if kid == "" {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(map[string]string{
			"msg": "kid can't be empty",
		})
	}
	arg := depositories.DepositoryCond{KID: kid}
	result, err := h.dbHandler.Get(arg)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		if err == pg.ErrNoRows {
			ctx.Status(http.StatusNotFound)
		}
		klog.Errorf("[Error] Get %s error %s", kid, err)
		return ctx.JSON(map[string]string{
			"msg": err.Error(),
		})
	}
	return ctx.JSON(result)
}
