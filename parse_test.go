package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseXML(t *testing.T) {
	f, err := ParseXML(nil)
	assert.Error(t, err)
	assert.Empty(t, f)
}

func TestParseXML2(t *testing.T) {
	r, err := os.Open("testdata/empty.txt")
	require.NoError(t, err)
	require.NotNil(t, r)

	// Prepare error
	require.NoError(t, os.Chmod("testdata/empty.txt", 0000))

	f, err := ParseXML(r)
	assert.Error(t, err)
	assert.Empty(t, f)

	require.NoError(t, os.Chmod("testdata/empty.txt", 0644))
}
