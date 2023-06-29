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
	"os"
	"time"

	"github.com/bestchains/bc-saas/pkg/models"
	"github.com/bestchains/bc-saas/pkg/utils"
	"github.com/go-pg/pg/v10"
	"k8s.io/klog/v2"
)

type dbHandler struct {
	// template images for depository in different style
	templateBaseImages map[Style]string
	// config for ttf font file
	ttfFontPath string

	db *pg.DB
}

func NewDBHandler(db *pg.DB, templateBaseImages map[Style]string, ttpfFontPath string) (Interface, error) {
	// check tempalte image and ttf font file exists
	for style, templateImage := range templateBaseImages {
		if _, ok := CertificateStyles[style]; !ok {
			return nil, fmt.Errorf("invalid certificate style %s", style)
		}
		if _, err := os.Stat(templateImage); os.IsNotExist(err) {
			return nil, fmt.Errorf("template image %s for %s does not exist", templateImage, style)
		}
	}
	if ttpfFontPath != "" {
		if _, err := os.Stat(ttpfFontPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("ttf font file %s does not exist", ttpfFontPath)
		}
	}

	return &dbHandler{db: db, templateBaseImages: templateBaseImages, ttfFontPath: ttpfFontPath}, nil
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

// GetCertificate get certificate for depository. Only support two styles: CN and ENG
func (h *dbHandler) GetCertificate(arg DepositoryCond, style Style) ([]byte, error) {
	// validate certificate style
	if style == "" {
		style = StyleCN
	}
	tpl, ok := CertificateStyles[style]
	if !ok {
		return nil, fmt.Errorf("invalid certificate style %s", style)
	}
	// generate certificate
	// FIXME: cache generated certificate
	template := &utils.GoPDFTemplate{}
	if err := template.Load([]byte(tpl)); err != nil {
		return nil, err
	}
	// set image
	template.Image = h.templateBaseImages[style]

	// get depository
	depository, err := h.Get(arg)
	if err != nil {
		return nil, err
	}

	return template.Render(utils.RenderOpts{
		TtfFontPath: h.ttfFontPath,
		Inputs: map[string]string{
			"name":             depository.Name,
			"owner":            depository.Owner,
			"kid":              depository.KID,
			"contentID":        depository.ContentID,
			"transactionHash":  depository.TransactionID,
			"trustedTimestamp": time.Unix(depository.TrustedTimestamp, 0).Format("2006-01-02 15:04:05"),
			"currentDate":      time.Now().Format("2006-01-02"),
		},
	})
}
