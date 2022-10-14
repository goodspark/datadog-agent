// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package helpers

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FlareBuilderMock struct {
	Fb   FlareBuilder
	Root string
}

func NewFlareBuilderMock(t *testing.T) FlareBuilder {
	root := t.TempDir()
	fb := &builder{
		fb:   newBuilder(root),
		root: root,
	}
	return fb
}

// AssertFileExists asserts that a file exists within the flare
func (mock *FlareBuilderMock) AssertFileExists(t *testing.T, path ...string) {
	assert.FileExists(t, filepath.Join(mock.Root, path...))
}

// AssertNoFileExists asserts that a file does not exists within the flare
func (mock *FlareBuilderMock) AssertNoFileExists(t *testing.T, path ...string) {
	assert.NoFileExists(t, filepath.Join(mock.Root, path...))
}
