package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/keltia/sandbox"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func Before(t *testing.T) *sandbox.Dir {
	snd, err := sandbox.New("test")
	require.NoError(t, err)

	return snd
}

func TestCheckFilename(t *testing.T) {
	td := []struct {
		In  string
		Out bool
	}{
		{"foo.bar", false},
		{"example.com!keltia.net!1538604008!1538690408.xml.gz", true},
	}
	for _, e := range td {
		res := checkFilename(e.In)
		assert.Equal(t, res, e.Out)
	}
}

func TestOpenFile(t *testing.T) {
	snd := Before(t)

	file := "/nonexistent"
	r, err := OpenFile(snd.Cwd(), file)
	assert.Error(t, err)
	assert.Nil(t, r)

	snd.Cleanup()
}

func TestOpenFile2(t *testing.T) {
	snd := Before(t)

	file := "testdata/empty.txt"
	file, err := filepath.Abs(file)
	require.NoError(t, err)

	err = snd.Enter()
	require.NoError(t, err)

	r, err := OpenFile(snd.Cwd(), file)
	defer r.Close()

	assert.NoError(t, err)
	assert.NotNil(t, r)

	snd.Cleanup()
}

func TestOpenFile3(t *testing.T) {
	snd := Before(t)

	file := "testdata/example.com!keltia.net!1538604008!1538690408.xml.gz"
	file, err := filepath.Abs(file)
	require.NoError(t, err)

	err = snd.Enter()
	require.NoError(t, err)

	r, err := OpenFile(snd.Cwd(), file)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	defer r.Close()

	bfile := filepath.Base(file)
	require.NotEmpty(t, bfile)

	ext := filepath.Ext(bfile)
	assert.Equal(t, ".gz", ext)

	pc := strings.Split(bfile, ".")
	assert.Equal(t, 5, len(pc))

	unc := strings.Join(pc[0:len(pc)-1], ".")

	assert.Equal(t, "example.com!keltia.net!1538604008!1538690408.xml", unc)
	assert.FileExists(t, filepath.Join(snd.Cwd(), unc))

	snd.Cleanup()
}

func TestHandleSingleFile(t *testing.T) {
	snd := Before(t)

	file := "empty.txt"
	txt, err := HandleSingleFile(snd, file)
	assert.Error(t, err)
	assert.Empty(t, txt)

	snd.Cleanup()
}

func TestOpenGzipFile(t *testing.T) {

}

func TestExtractXML(t *testing.T) {

}
