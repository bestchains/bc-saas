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
	"testing"

	"github.com/signintech/gopdf"
	"github.com/stretchr/testify/assert"
)

// TestGoPDFTemplate_Load tests the Load method of GoPDFTemplate
func TestGoPDFTemplate_Load(t *testing.T) {
	// Arrange
	gop := &GoPDFTemplate{}
	bytes := []byte(`{
    "image": "./testdata/template.jpg",
    "locations": [
        {
            "text": "Hello %s",
            "inputs": [
                "world"
            ],
            "x": 10,
            "y": 20,
            "style": "",
            "size": 12
        }
    ]
}`)

	// Act
	err := gop.Load(bytes)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "./testdata/template.jpg", gop.Image)
	assert.Len(t, gop.Locations, 1)
	assert.Equal(t, "Hello %s", gop.Locations[0].Text)
	assert.Equal(t, []string{"world"}, gop.Locations[0].Inputs)
	assert.Equal(t, float64(10), gop.Locations[0].X)
	assert.Equal(t, float64(20), gop.Locations[0].Y)
	assert.Equal(t, 12, gop.Locations[0].Size)
}

// TestGoPDFTemplate_LoadFromFile tests the LoadFromFile method of GoPDFTemplate
func TestGoPDFTemplate_LoadFromFile(t *testing.T) {
	// Arrange
	gop := &GoPDFTemplate{}
	path := "./testdata/template.json"

	// Act
	err := gop.LoadFromFile(path)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "./testdata/template.jpg", gop.Image)
	assert.Len(t, gop.Locations, 1)
	assert.Equal(t, "Hello %s", gop.Locations[0].Text)
	assert.Equal(t, []string{"world"}, gop.Locations[0].Inputs)
	assert.Equal(t, float64(10), gop.Locations[0].X)
	assert.Equal(t, float64(20), gop.Locations[0].Y)
	assert.Equal(t, 12, gop.Locations[0].Size)
}

// TestGoPDFTemplate_Render tests the Render method of GoPDFTemplate
func TestGoPDFTemplate_Render(t *testing.T) {
	// Arrange
	gop := &GoPDFTemplate{
		Image: "./testdata/template.jpg",
		Locations: []Location{
			{
				Text:   "Hello %s",
				Inputs: []string{"name"},
				X:      10,
				Y:      20,
				Size:   12,
			},
		},
	}
	option := RenderOpts{
		PageSize:    *gopdf.PageSizeA4,
		TtfFontPath: "./ttf/stsong.ttf",
		Inputs: map[string]string{
			"name": "world",
		},
	}

	// Act
	bytes, err := gop.Render(option)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, bytes)
}
