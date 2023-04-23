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

package models

// Depository defines valuable fields for a depository
type Depository struct {
	Index       string `json:"index" pg:"index"`
	KID         string `json:"kid" pg:"kid,pk"`
	Platform    string `json:"platform" pg:"platform"`
	Operator    string `json:"operator" pg:"operator"`
	Owner       string `json:"owner" pg:"owner"`
	BlockNumber string `json:"blockNumber" pg:"blockNumber"`

	// Content related
	ContentName      string `json:"contentName" pg:"contentName"`
	ContentID        string `json:"contentID" pg:"contentID"`
	ContentType      string `json:"contentType" pg:"contentType"`
	TrustedTimestamp string `json:"trustedTimestamp" pg:"trustedTimestamp"`
}
