package main

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyze(t *testing.T) {
	s, err := Analyze(Feedback{})
	assert.Error(t, err)
	assert.Empty(t, s)
}

func TestGatherRows_Empty(t *testing.T) {
	r := GatherRows(Feedback{})
	assert.Empty(t, r)
}

func TestGatherRows_Good(t *testing.T) {
	file := "testdata/example.com!keltia.net!1538604008!1538690408.xml"

	a, err := NewArchive(file)
	require.NoError(t, err)

	body, err := a.Extract(".xml")
	require.NoError(t, err)

	var report Feedback

	err = xml.Unmarshal(body, &report)
	require.NoError(t, err)

	rows := GatherRows(report)
	assert.Equal(t, 1, len(rows))
}
