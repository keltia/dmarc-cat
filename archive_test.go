package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewArchive_Empty(t *testing.T) {
	a, err := NewArchive("")
	require.Error(t, err)
	assert.Empty(t, a)
	assert.IsType(t, (*Plain)(nil), a)
}

func TestNewArchive_Plain(t *testing.T) {
	a, err := NewArchive("foo.txt")
	require.NoError(t, err)
	assert.NotEmpty(t, a)
	assert.IsType(t, (*Plain)(nil), a)
}

func TestNewArchive_Zip(t *testing.T) {
	a, err := NewArchive("testdata/google.com!keltia.net!1538438400!1538524799.zip")
	require.NoError(t, err)
	assert.NotEmpty(t, a)
	assert.IsType(t, (*Zip)(nil), a)

}

func TestNewArchive_ZipNone(t *testing.T) {
	a, err := NewArchive("foo.zip")
	require.Error(t, err)
	assert.Empty(t, a)
	assert.IsType(t, (*Zip)(nil), a)

}

func TestNewArchive_Gzip(t *testing.T) {
	a, err := NewArchive("foo.gz")
	require.NoError(t, err)
	assert.NotEmpty(t, a)
	assert.IsType(t, (*Gzip)(nil), a)

}

// Plain

func TestPlain_Extract(t *testing.T) {
	fn := "testdata/notempty.txt"
	a, err := NewArchive(fn)
	require.NoError(t, err)
	require.NotNil(t, a)

	txt, err := a.Extract("")
	assert.NoError(t, err)
	assert.Equal(t, "this is a file\n", string(txt))
}

func TestPlain_Extract2(t *testing.T) {
	fn := "testdata/notempty.txt"
	a, err := NewArchive(fn)
	require.NoError(t, err)
	require.NotNil(t, a)

	txt, err := a.Extract(".doc")
	assert.Error(t, err)
	assert.Empty(t, txt)
}

func TestPlain_Close(t *testing.T) {
	fn := "testdata/notempty.txt"
	a, err := NewArchive(fn)
	require.NoError(t, err)
	require.NotNil(t, a)

	require.NoError(t, a.Close())
}

// Zip

func TestZip_Extract(t *testing.T) {
	fn := "testdata/google.com!keltia.net!1538438400!1538524799.zip"
	a, err := NewArchive(fn)
	require.NoError(t, err)
	require.NotNil(t, a)
	defer a.Close()

	rh, err := ioutil.ReadFile("testdata/google.com!keltia.net!1538438400!1538524799.xml")
	require.NoError(t, err)
	require.NotEmpty(t, rh)

	txt, err := a.Extract(".xml")
	assert.NoError(t, err)
	assert.Equal(t, string(rh), string(txt))
}

func TestZip_Extract2(t *testing.T) {
	fn := "testdata/notempty.zip"
	a, err := NewArchive(fn)
	require.NoError(t, err)
	require.NotNil(t, a)

	txt, err := a.Extract(".xml")
	assert.Error(t, err)
	assert.Empty(t, txt)
}

func TestZip_Extract3(t *testing.T) {
	fn := "testdata/google.com!keltia.net!1538438400!1538524799.zip"
	a, err := NewArchive(fn)
	require.NoError(t, err)
	require.NotNil(t, a)
	defer a.Close()

	rh, err := ioutil.ReadFile("testdata/google.com!keltia.net!1538438400!1538524799.xml")
	require.NoError(t, err)
	require.NotEmpty(t, rh)

	txt, err := a.Extract(".txt")
	assert.Error(t, err)
	assert.Empty(t, txt)
}

func TestZip_Close(t *testing.T) {
	fn := "testdata/notempty.zip"
	a, err := NewArchive(fn)
	require.NoError(t, err)
	require.NotNil(t, a)

	require.NoError(t, a.Close())
}

// Gzip

func TestGzip_Extract(t *testing.T) {
	fn := "testdata/example.com!keltia.net!1538604008!1538690408.xml.gz"
	a, err := NewArchive(fn)
	require.NoError(t, err)
	require.NotNil(t, a)
	defer a.Close()

	rh, err := ioutil.ReadFile("testdata/example.com!keltia.net!1538604008!1538690408.xml")
	require.NoError(t, err)
	require.NotEmpty(t, rh)

	txt, err := a.Extract(".xml")
	assert.NoError(t, err)
	assert.Equal(t, string(rh), string(txt))
}

func TestGzip_Extract2(t *testing.T) {
	fn := "testdata/notempty.txt.gz"
	a, err := NewArchive(fn)
	require.NoError(t, err)
	require.NotNil(t, a)

	rh, err := ioutil.ReadFile("testdata/notempty.txt")
	require.NoError(t, err)
	require.NotEmpty(t, rh)

	txt, err := a.Extract(".txt")
	assert.NoError(t, err)
	assert.NotEmpty(t, txt)
	assert.Equal(t, string(rh), string(txt))
}

func TestGzip_Extract3(t *testing.T) {
	fn := "/nonexistent"
	a, err := NewArchive(fn)
	require.NoError(t, err)
	require.NotNil(t, a)

	txt, err := a.Extract(".txt")
	assert.Error(t, err)
	assert.Empty(t, txt)
}

func TestGzip_Close(t *testing.T) {
	fn := "testdata/notempty.txt.gz"
	a, err := NewArchive(fn)
	require.NoError(t, err)
	require.NotNil(t, a)

	require.NoError(t, a.Close())
}
