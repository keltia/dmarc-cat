package main

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyze(t *testing.T) {
	ctx := &Context{NullResolver{}, nil}
	s, err := Analyze(ctx, Feedback{})
	assert.Error(t, err)
	assert.Empty(t, s)
}

func TestGatherRows_Empty(t *testing.T) {
	ctx := &Context{NullResolver{}, nil}
	r := GatherRows(ctx, Feedback{})
	assert.Empty(t, r)
}

func TestGatherRows_Good(t *testing.T) {
	ctx := &Context{NullResolver{}, nil}
	file := "testdata/example.com!keltia.net!1538604008!1538690408.xml"

	a, err := NewArchive(file)
	require.NoError(t, err)

	body, err := a.Extract(".xml")
	require.NoError(t, err)

	var report Feedback

	err = xml.Unmarshal(body, &report)
	require.NoError(t, err)

	rows := GatherRows(ctx, report)
	assert.Equal(t, 1, len(rows))
}

type ErrResolver struct{}

func (ErrResolver) LookupAddr(ip string) ([]string, error) {
	return []string{}, fmt.Errorf("fake error")
}

func TestResolveIP_Error(t *testing.T) {
	ctx := &Context{r: ErrResolver{}}
	ips := ResolveIP(ctx, "8.8.8.8")
	assert.Equal(t, "8.8.8.8", ips)
}

func TestResolveIP_Good(t *testing.T) {
	ctx := &Context{r: ErrResolver{}}
	ips := ResolveIP(ctx, "8.8.8.8")
	assert.NotEmpty(t, ips)
	assert.Equal(t, "8.8.8.8", ips)
}
