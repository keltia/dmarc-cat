package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyze(t *testing.T) {

}

func TestGatherRows_Empty(t *testing.T) {
	r := GatherRows(Feedback{})
	assert.Empty(t, r)
}
