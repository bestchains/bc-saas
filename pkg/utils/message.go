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

package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

var (
	ErrNotMessage            = errors.New("not a message")
	ErrInvalidMessage        = errors.New("invalid message")
	ErrInvalidMessageSender  = errors.New("invalid message sender")
	ErrAlgorithmNotSupported = errors.New("algorithm not supported yet")
	ErrInvalidSignature      = errors.New("invalid signature")
)

type Message struct {
	Nonce     uint64 `json:"nonce"`
	PublicKey []byte `json:"publicKey"`
	Signature []byte `json:"signature"`
}

func (msg *Message) Marshal() ([]byte, error) {
	if msg == nil {
		msg = new(Message)
	}
	return json.Marshal(msg)
}

func (msg *Message) Unmarshal(bytes []byte) error {
	var err error
	if msg == nil {
		msg = new(Message)
	}
	if err = json.Unmarshal(bytes, msg); err != nil {
		return errors.Wrap(ErrNotMessage, err.Error())
	}

	return nil
}

func (msg *Message) VerifyAgainstArgs(args ...string) (string, error) {
	payload := msg.GeneratePayload(args...)

	pub, err := x509.ParsePKIXPublicKey(msg.PublicKey)
	if err != nil {
		return "", errors.Wrap(ErrInvalidMessage, err.Error())
	}

	msgSender, err := FromPublicKey(pub)
	if err != nil {
		return "", err
	}

	switch pub := pub.(type) {
	case *ecdsa.PublicKey:
		hashedPayload := GenerateHash(payload)
		if !ecdsa.VerifyASN1(pub, hashedPayload, msg.Signature) {
			return "", errors.Wrap(ErrInvalidMessage, ErrInvalidSignature.Error())
		}
	default:
		return "", ErrAlgorithmNotSupported
	}

	return msgSender, nil
}
func (msg *Message) GeneratePayload(args ...string) []byte {
	payload := []byte(strconv.FormatUint(msg.Nonce, 10))
	for _, arg := range args {
		payload = append(payload, []byte(arg)...)
	}
	return payload
}

func GenerateHash(payload []byte) []byte {
	return sha512.New().Sum(payload[:])
}

func FromPublicKey(pub interface{}) (string, error) {
	publicKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return "", ErrAlgorithmNotSupported
	}
	// Serialize the public key
	serializedPubKey := elliptic.Marshal(elliptic.P256(), publicKey.X, publicKey.Y)

	// Hash the public key using Keccak-256
	hashedPubKey := sha3.Sum256(serializedPubKey)

	// Truncate the hash and add a prefix to get the final Ethereum address
	address := "0x" + hex.EncodeToString(hashedPubKey[12:])

	return address, nil
}
