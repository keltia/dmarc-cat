package main

import (
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

func TestNewArchive_Tar(t *testing.T) {
	a, err := NewArchive("foo.tar")
	require.NoError(t, err)
	assert.NotEmpty(t, a)
	assert.IsType(t, (*Tar)(nil), a)

}

func TestPlain_Extract(t *testing.T) {
	fn := "testdata/notempty.txt"
	a, err := NewArchive(fn)
	require.NoError(t, err)
	require.NotNil(t, a)

	txt, err := a.Extract("")
	assert.Equal(t, "this is a file\n", string(txt))
}
