package parser

import (
	"testing"

	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

func TestParseGoMod(t *testing.T) {
	path := "go.mod"
	data := []byte(`module github.com/harryzcy/test-module

	go 1.19
	
	require (
		github.com/docker/docker v20.10.11+incompatible
		github.com/spf13/viper v1.14.0
		github.com/stretchr/testify v1.8.1
	)
	
	require (
		github.com/davecgh/go-spew v1.1.1 // indirect
	)`)
	want := []types.Dependency{
		{
			File:           "go.mod",
			Name:           "github.com/docker/docker",
			Version:        "v20.10.11+incompatible",
			Indirect:       false,
			PackageManager: "go-mod",
			Position: types.Position{
				Line:      6,
				StartByte: 64,
				EndByte:   111,
			},
		},
		{
			File:           "go.mod",
			Name:           "github.com/spf13/viper",
			Version:        "v1.14.0",
			Indirect:       false,
			PackageManager: "go-mod",
			Position: types.Position{
				Line:      7,
				StartByte: 114,
				EndByte:   144,
			},
		},
		{
			File:           "go.mod",
			Name:           "github.com/stretchr/testify",
			Version:        "v1.8.1",
			Indirect:       false,
			PackageManager: "go-mod",
			Position: types.Position{
				Line:      8,
				StartByte: 147,
				EndByte:   181,
			},
		},
		{
			File:           "go.mod",
			Name:           "github.com/davecgh/go-spew",
			Version:        "v1.1.1",
			Indirect:       true,
			PackageManager: "go-mod",
			Position: types.Position{
				Line:      12,
				StartByte: 200,
				EndByte:   233,
			},
		},
	}

	dependencies, err := parseGoMod(path, data)
	assert.Nil(t, err)
	assert.Equal(t, want, dependencies)
}
