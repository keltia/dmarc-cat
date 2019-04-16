package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/keltia/archive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestHandleZipFile(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}

	file := "testdata/google.com!keltia.net!1538438400!1538524799.zip"

	txt, err := HandleZipFile(ctx, file)
	assert.NoError(t, err)
	assert.NotEmpty(t, txt)
}

func TestHandleZipFile_Xml(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}

	file := "testdata/example.com!keltia.net!1538604008!1538690408.xml"

	txt, err := HandleZipFile(ctx, file)
	assert.NoError(t, err)
	assert.NotEmpty(t, txt)
}

func TestHandleZipFile_Bad(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}

	file := "testdata/notempty.zip"

	txt, err := HandleZipFile(ctx, file)
	assert.Error(t, err)
	assert.Empty(t, txt)
}

func TestHandleZipFile_Bad1(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}

	file := "testdata/bad.zip"

	txt, err := HandleZipFile(ctx, file)
	assert.Error(t, err)
	assert.Empty(t, txt)
}

func TestHandleZipFile_None(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}

	file := "/nonexistent"

	txt, err := HandleZipFile(ctx, file)
	assert.Error(t, err)
	assert.Empty(t, txt)
}

func TestHandleSingleFile_Plain(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}

	file := "testdata/empty.txt"
	fh, err := os.Open(file)
	require.NoError(t, err)
	txt, err := HandleSingleFile(ctx, fh, archive.Ext2Type(filepath.Ext(file)))
	assert.Error(t, err)
	assert.Empty(t, txt)
}

func TestHandleSingleFile_Gzip(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}

	file := "testdata/example.com!keltia.net!1538604008!1538690408.xml.gz"

	fh, err := os.Open(file)
	require.NoError(t, err)
	txt, err := HandleSingleFile(ctx, fh, archive.Ext2Type(".gz"))
	assert.NoError(t, err)
	assert.NotEmpty(t, txt)
}

func TestHandleSingleFile_Zip(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}

	file := "testdata/google.com!keltia.net!1538438400!1538524799.zip"

	fh, err := os.Open(file)
	require.NoError(t, err)
	txt, err := HandleSingleFile(ctx, fh, archive.Ext2Type(filepath.Ext(file)))
	assert.Error(t, err)
	assert.Empty(t, txt)
}

func TestHandleSingleFile_Xml(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}

	fDebug = true
	file := "testdata/example.com!keltia.net!1538604008!1538690408.xml"

	fh, err := os.Open(file)
	require.NoError(t, err)

	assert.Equal(t, archive.ArchivePlain, archive.Ext2Type(filepath.Ext(file)))

	txt, err := HandleSingleFile(ctx, fh, archive.Ext2Type(filepath.Ext(file)))
	assert.NoError(t, err)
	assert.NotEmpty(t, txt)
	fDebug = false
}

func TestHandleSingleFile_Verbose(t *testing.T) {
	fVerbose = true

	ctx := &Context{NullResolver{}, 1}
	file := "testdata/empty.txt"

	fh, err := os.Open(file)
	require.NoError(t, err)

	txt, err := HandleSingleFile(ctx, fh, archive.Ext2Type(filepath.Ext(file)))
	assert.Error(t, err)
	assert.Empty(t, txt)

	fVerbose = false
}

func TestHandleSingleFile_Debug(t *testing.T) {
	fDebug = true

	ctx := &Context{NullResolver{}, 1}

	file := "testdata/empty.txt"

	fh, err := os.Open(file)
	require.NoError(t, err)

	txt, err := HandleSingleFile(ctx, fh, archive.Ext2Type(filepath.Ext(file)))
	assert.Error(t, err)
	assert.Empty(t, txt)

	fDebug = false
}
