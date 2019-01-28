package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFilename(t *testing.T) {
	td := []struct {
		In  string
		Out bool
	}{
		{"foo.bar", false},
		{"example.com!keltia.net!1538604008!1538690408.xml.gz", true},
		{"example.com!keltia.net!1538604008!1538690408.xml.gz", true},
		{"example.com!keltia.net!1538604008!1538690408.xml", true},
		{"google.com!keltia.net!1538438400!1538524799.zip", true},
	}
	for _, e := range td {
		res := checkFilename(e.In)
		assert.Equal(t, res, e.Out)
	}
}

func TestHandleSingleFile(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}

	file := "empty.txt"
	txt, err := HandleSingleFile(ctx, file)
	assert.Error(t, err)
	assert.Empty(t, txt)
}

func TestHandleSingleFile2(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}

	file := "testdata/example.com!keltia.net!1538604008!1538690408.xml.gz"
	txt, err := HandleSingleFile(ctx, file)
	assert.NoError(t, err)
	assert.NotEmpty(t, txt)
}

func TestHandleSingleFile3(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}

	file := "testdata/google.com!keltia.net!1538438400!1538524799.zip"
	txt, err := HandleSingleFile(ctx, file)
	assert.NoError(t, err)
	assert.NotEmpty(t, txt)
}

func TestHandleSingleFile4(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}

	file := "testdata/example.com!keltia.net!1538604008!1538690408.xml"
	txt, err := HandleSingleFile(ctx, file)
	assert.NoError(t, err)
	assert.NotEmpty(t, txt)
}

func TestHandleSingleFile_Verbose(t *testing.T) {
	fVerbose = true

	ctx := &Context{NullResolver{}, 1}
	file := "empty.txt"
	txt, err := HandleSingleFile(ctx, file)
	assert.Error(t, err)
	assert.Empty(t, txt)

	fVerbose = false
}

func TestHandleSingleFile_Debug(t *testing.T) {
	fDebug = true

	ctx := &Context{NullResolver{}, 1}

	file := "empty.txt"
	txt, err := HandleSingleFile(ctx, file)
	assert.Error(t, err)
	assert.Empty(t, txt)

	fDebug = false
}
