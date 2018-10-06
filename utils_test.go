package main

import (
	"testing"
)

func TestVerbose_No(t *testing.T) {
	verbose("no")
}

func TestVerbose_Yes(t *testing.T) {
	fVerbose = true
	verbose("yes")
	fVerbose = false
}

func TestDebug_No(t *testing.T) {
	debug("no")
}

func TestDebug_Yes(t *testing.T) {
	fVerbose = true
	verbose("yes")
	fVerbose = false
}
