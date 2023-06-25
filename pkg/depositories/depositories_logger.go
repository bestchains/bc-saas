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

import "github.com/bestchains/bc-saas/pkg/models"

type loggerHandler struct{}

func NewLoggerHandler() Interface {
	return &loggerHandler{}
}

func (l *loggerHandler) List(arg DepositoryCond) ([]models.Depository, int64, error) {
	return nil, 0, nil
}
func (l *loggerHandler) Get(arg DepositoryCond) (models.Depository, error) {
	return models.Depository{}, nil
}

func (l *loggerHandler) GetCertificate(arg DepositoryCond, language string) ([]byte, error) {
	return []byte{}, nil
}
