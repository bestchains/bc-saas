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
	"github.com/bestchains/bc-saas/pkg/models"
	"github.com/go-pg/pg/v10"
	"k8s.io/klog/v2"
)

type dbHandler struct {
	db *pg.DB
}

func NewDBHandler(db *pg.DB) Interface {
	return &dbHandler{db: db}
}

func (h *dbHandler) List(arg DepositoryCond) ([]models.Depository, int64, error) {
	result := make([]models.Depository, 0)
	cond, params := arg.ToCond()
	klog.V(5).Infof(" dbHandler list query %v %v\n", cond, params)

	q := h.db.Model(&result)
	for i := 0; i < len(cond); i++ {
		q = q.Where(cond[i], params[i])
	}
	c, err := q.Count()
	if err != nil {
		return result, 0, err
	}
	q = q.Order(`trustedTimestamp desc`)
	if arg.Size != 0 {
		q = q.Limit(arg.Size).Offset(arg.From)
	}
	if err := q.Select(); err != nil {
		return result, 0, err
	}

	return result, int64(c), nil
}

func (h *dbHandler) Get(arg DepositoryCond) (models.Depository, error) {
	result := models.Depository{}
	cond, params := arg.ToCond()
	q := h.db.Model(&result)
	for i := 0; i < len(cond); i++ {
		q = q.Where(cond[i], params[i])
	}
	err := q.Select()
	return result, err
}
