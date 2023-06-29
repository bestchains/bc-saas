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
	"encoding/json"
	"fmt"
	"os"

	"github.com/signintech/gopdf"
)

// GoPDFTemplate is a template for gopdf to render a pdf from a image
type GoPDFTemplate struct {
	// Path to image template
	Image string `json:"image,omitempty"`
	// Locations to write text
	Locations []Location `json:"locations,omitempty"`
}

// Location is a location to write text
type Location struct {
	// Text to be written
	Text string `json:"text"`
	// Inputs to be parsed into text with fmt.Sprintf()
	Inputs []string `json:"inputs"`
	// Position of text
	X float64 `json:"x"`
	Y float64 `json:"y"`
	// Font style and size
	Style string `json:"style"`
	Size  int    `json:"size"`
	// CellSize
	CellSize *gopdf.Rect `json:"cell_size,omitempty"`
}

// Load loads a GoPDFTemplate from bytes
func (gop *GoPDFTemplate) Load(bytes []byte) error {
	err := json.Unmarshal(bytes, gop)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	return nil
}

// LoadFromFile loads a GoPDFTemplate from a file
func (gop *GoPDFTemplate) LoadFromFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, gop)
	if err != nil {
		return err
	}
	return nil
}

var (
	defaultPageSize = *gopdf.PageSizeA4
	defaultCellSize = gopdf.Rect{
		W: 420,
		H: 100,
	}
)

const (
	// defaultTtfFontPath is the default path to ttf font
	// stsong.ttf downloaded from https://www.wfonts.com/font/stsong
	// Reason to use this font: it supports Chinese(https://github.com/signintech/gopdf/issues/68)
	defaultTtfFontPath = "resource/ttf/stsong.ttf"
)

// RenderOpts is the options for GoPDFTemplate.Render()
type RenderOpts struct {
	PageSize gopdf.Rect
	// Path to ttf font
	TtfFontPath string
	// Inputs to be parsed into text with fmt.Sprintf()
	Inputs map[string]string
	// Cell size for gopdf to render the pdf
	CellSize gopdf.Rect
}

// Render renders a pdf from a GoPDFTemplate
func (gop *GoPDFTemplate) Render(option RenderOpts) ([]byte, error) {
	var err error

	if option.PageSize == (gopdf.Rect{}) {
		option.PageSize = defaultPageSize
	}
	if option.CellSize == (gopdf.Rect{}) {
		option.CellSize = defaultCellSize
	}
	if option.TtfFontPath == "" {
		option.TtfFontPath = defaultTtfFontPath
	}
	if option.Inputs == nil {
		option.Inputs = make(map[string]string)
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: option.PageSize})
	pdf.AddPage()

	// load image template
	// if image is not set, skip loading image
	if gop.Image != "" {
		pdf.Image(gop.Image, 0, 0, gopdf.PageSizeA4)
	}

	// load ttf font
	if err = pdf.AddTTFFontWithOption("font", option.TtfFontPath, gopdf.TtfOption{UseKerning: true}); err != nil {
		return nil, fmt.Errorf("failed to load ttf font: %w", err)
	}

	// set words based on locations
	for _, location := range gop.Locations {
		err = pdf.SetFont("font", location.Style, location.Size)
		if err != nil {
			return nil, fmt.Errorf("failed to set font: %w", err)
		}

		// set position
		pdf.SetXY(location.X, location.Y)

		// parse text with args
		inputs := make([]any, len(location.Inputs))
		for i, arg := range location.Inputs {
			inputs[i] = option.Inputs[arg]
		}
		text := fmt.Sprintf(location.Text, inputs...)

		// set cell size
		cellSize := option.CellSize
		if location.CellSize != nil {
			cellSize = *location.CellSize
		}

		// set text with limited witdh and height
		// all cells has same cell size
		if err = pdf.MultiCellWithOption(&cellSize, text, gopdf.CellOption{}); err != nil {
			return nil, fmt.Errorf("failed to set text: %w", err)
		}
	}

	return pdf.GetBytesPdf(), nil
}
