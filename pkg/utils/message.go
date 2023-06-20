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

// Package utils provides utility functions for working with cryptographic messages
package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

var (
	// ErrNotMessage is returned when trying to unmarshal invalid message bytes
	ErrNotMessage = errors.New("not a message")
	// ErrInvalidMessage is returned when trying to verify an invalid message
	ErrInvalidMessage = errors.New("invalid message")
	// ErrInvalidMessageSender is returned when trying to verify a message with an invalid sender
	ErrInvalidMessageSender = errors.New("invalid message sender")
	// ErrAlgorithmNotSupported is returned when trying to use an unsupported cryptographic algorithm
	ErrAlgorithmNotSupported = errors.New("algorithm not supported yet")
	// ErrInvalidSignature is returned when trying to verify a message with an invalid signature
	ErrInvalidSignature = errors.New("invalid signature")
)

// Message represents a cryptographic message
type Message struct {
	Nonce     uint64 `json:"nonce"`
	PublicKey []byte `json:"publicKey"`
	Signature []byte `json:"signature"`
}

// Marshal returns the JSON encoding of the message.
// If the message is nil, a new message is created.
func (msg *Message) Marshal() ([]byte, error) {
	// Marshal the message as JSON
	return json.Marshal(msg)
}

// Unmarshal unmarshals the given bytes into a Message struct.
// It returns an error if the unmarshalling fails.
func (msg *Message) Unmarshal(bytes []byte) error {
	var err error
	if err = json.Unmarshal(bytes, msg); err != nil {
		return errors.Wrap(ErrNotMessage, err.Error())
	}

	return nil
}

// UnmarshalBase64Str unmarshals the given base64-encoded string into a Message struct.
// It returns an error if the unmarshalling fails.
func (msg *Message) UnmarshalBase64Str(str string) error {
	// base64 decode
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return errors.Wrap(ErrNotMessage, err.Error())
	}
	return msg.Unmarshal(bytes)
}

// VerifyAgainstArgs verifies that the message signature is valid against the given arguments.
// It returns the Ethereum address of the message sender and an error, if any.
func (msg *Message) VerifyAgainstArgs(args ...string) (string, error) {
	payload := msg.GeneratePayload(args...)

	// Parse the public key
	pub, err := x509.ParsePKIXPublicKey(msg.PublicKey)
	if err != nil {
		return "", errors.Wrap(ErrInvalidMessage, err.Error())
	}

	// Get the Ethereum address of the message sender
	msgSender, err := FromPublicKey(pub)
	if err != nil {
		return "", err
	}

	// Verify the message signature
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

// GeneratePayload generates a payload for the message with the given arguments.
// The payload includes the message nonce and all the arguments appended together.
func (msg *Message) GeneratePayload(args ...string) []byte {
	payload := []byte(strconv.FormatUint(msg.Nonce, 10))
	for _, arg := range args {
		payload = append(payload, []byte(arg)...)
	}
	return payload
}

// GenerateHash generates a SHA512 hash for the given payload.
func GenerateHash(payload []byte) []byte {
	return sha512.New().Sum(payload[:])
}

// FromPublicKey generates an Ethereum address from the given public key.
// It serializes the public key, hashes it using Keccak-256, truncates the hash, and adds a prefix to get the final Ethereum address.
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
