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

import (
	"context"

	"github.com/go-pg/pg/v10"
	"k8s.io/klog/v2"
)

// Depository defines valuable fields for a depository
type Depository struct {
	Index       string `json:"index" pg:"index"`
	KID         string `json:"kid" pg:"kid,pk"`
	Platform    string `json:"platform" pg:"platform"`
	Operator    string `json:"operator" pg:"operator"`
	Owner       string `json:"owner" pg:"owner"`
	BlockNumber uint64 `json:"blockNumber" pg:"blockNumber"`

	// Content related
	Name             string `json:"name" pg:"name"`
	ContentName      string `json:"contentName" pg:"contentName"`
	ContentID        string `json:"contentID" pg:"contentID"`
	ContentType      string `json:"contentType" pg:"contentType"`
	TrustedTimestamp string `json:"trustedTimestamp" pg:"trustedTimestamp"`
	Description      string `json:"description" pg:"description"`
}

var _ pg.QueryHook = (*Depository)(nil)

func (*Depository) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	query, err := event.FormattedQuery()
	if err != nil {
		return ctx, nil
	}
	klog.V(5).Infof("[format query] %s\n", string(query))
	return ctx, nil
}

func (*Depository) AfterQuery(context.Context, *pg.QueryEvent) error {
	return nil
}
func MaxBlockNumber(db *pg.DB) uint64 {
	ans := uint64(0)
	if err := db.Model((*Depository)(nil)).ColumnExpr(`max("blockNumber") as bn`).Select(&ans); err != nil {
		klog.Errorf("[Error] select max blockNumber failed. error %s", err.Error())
	}
	return ans
}
