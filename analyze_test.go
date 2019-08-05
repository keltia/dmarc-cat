package main

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/keltia/archive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyze(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}
	s, err := Analyze(ctx, Feedback{})
	assert.Error(t, err)
	assert.Empty(t, s)
}

func TestGatherRows_Empty(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}
	r := GatherRows(ctx, Feedback{})
	assert.Empty(t, r)
}

func TestGatherRows_Good(t *testing.T) {
	ctx := &Context{NullResolver{}, 1}
	file := "testdata/example.com!keltia.net!1538604008!1538690408.xml"

	a, err := archive.New(file)
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
	return []string{"BAD"}, fmt.Errorf("fake error")
}

func TestParallelSolve_Error(t *testing.T) {
	ctx := &Context{r: ErrResolver{}, jobs: 1}

	td := []IP{
		{IP: "8.8.8.8", Name: "BAD"},
		{IP: "8.8.4.4", Name: "BAD"},
	}
	ips := ParallelSolve(ctx, td)
	assert.NotEmpty(t, ips)
	assert.EqualValues(t, td, ips)
}

func TestParallelSolve_Good(t *testing.T) {
	ctx := &Context{r: FakeResolver{}, jobs: 1}

	td := []IP{
		{IP: "8.8.8.8", Name: "8.8.8.8"},
		{IP: "8.8.4.4", Name: "8.8.4.4"},
	}
	ips := ParallelSolve(ctx, td)
	assert.NotEmpty(t, ips)
	assert.EqualValues(t, td, ips)
}
