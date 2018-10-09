package main

import (
	"os"
	"testing"
)

// XXX I'm not sure how to really test main() — os.Args() is global and not reset between calls

func TestMain_Noargs(t *testing.T) {
	main()
}

func TestMain_Noargs_Verbose(t *testing.T) {
	fVerbose = true
	main()
	fVerbose = false
}

func TestMain_Noargs_Debug(t *testing.T) {
	fDebug = true
	main()
	fDebug = false
}

func TestMain_Noargs_NoResolv(t *testing.T) {
	os.Args = append(os.Args, "testdata/google.com!keltia.net!1538438400!1538524799.zip")

	fNoResolv = true
	main()
	fNoResolv = false
}

func TestMain_GoodFile(t *testing.T) {
	os.Args = append(os.Args, "testdata/google.com!keltia.net!1538438400!1538524799.zip")
	main()
}

func TestMain_NoFile(t *testing.T) {
	os.Args = append(os.Args, "/nonexistent")
	main()
}

func TestMain_EmptyArg(t *testing.T) {
	os.Args = append(os.Args, "foo")
	main()
}
