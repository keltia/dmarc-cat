package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
