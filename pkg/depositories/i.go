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

package depositories

import (
	"fmt"

	"github.com/bestchains/bc-saas/pkg/models"
)

type DepositoryCond struct {
	From, Size             int
	Name, KID, ContentName string
	StartTime, EndTime     int64
}

func (dc *DepositoryCond) ToCond() ([]string, []interface{}) {
	params := make([]interface{}, 0)
	cond := make([]string, 0)
	if dc.KID != "" {
		cond = append(cond, `kid=?`)
		params = append(params, dc.KID)
	}
	if dc.ContentName != "" {
		cond = append(cond, `"contentName" like ?`)
		params = append(params, fmt.Sprintf(`%%%s%%`, dc.ContentName))
	}

	if dc.StartTime > 0 {
		cond = append(cond, `"trustedTimestamp">=?`)
		params = append(params, dc.StartTime)
	}
	if dc.EndTime > 0 {
		cond = append(cond, `"trustedTimestamp"<=?`)
		params = append(params, dc.EndTime)
	}

	if len(dc.Name) > 0 {
		cond = append(cond, `"name" like ?`)
		params = append(params, fmt.Sprintf(`%%%s%%`, dc.Name))
	}
	return cond, params
}

type Interface interface {
	Get(DepositoryCond) (models.Depository, error)
	GetCertificate(cond DepositoryCond, style Style) ([]byte, error)
	List(DepositoryCond) ([]models.Depository, int64, error)
}
