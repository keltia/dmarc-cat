package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullResolver_LookupAddr(t *testing.T) {
	var r NullResolver

	resp, err := r.LookupAddr("example.com")
	assert.NoError(t, err)
	assert.Equal(t, []string{"example.com"}, resp)
}

func TestRealResolver_LookupAddr(t *testing.T) {
	var r RealResolver

	resp, err := r.LookupAddr("8.8.8.8")
	assert.NoError(t, err)
	assert.Equal(t, []string{"dns.google."}, resp)
}
