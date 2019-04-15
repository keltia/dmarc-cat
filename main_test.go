package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	ctx, err := Setup([]string{})
	assert.Nil(t, ctx)
	assert.Error(t, err)
}

func TestSetup2(t *testing.T) {
	ctx, err := Setup([]string{"foo.zip"})
	assert.NotNil(t, ctx)
	assert.NoError(t, err)
	assert.IsType(t, (*Context)(nil), ctx)
}

func TestSetup3(t *testing.T) {
	fNoResolv = true
	ctx, err := Setup([]string{"foo.zip"})
	assert.NotNil(t, ctx)
	assert.NoError(t, err)
	assert.IsType(t, (*Context)(nil), ctx)
	fNoResolv = false
}

func TestSetup4(t *testing.T) {
	fVersion = true
	ctx, err := Setup([]string{"foo.zip"})
	assert.Nil(t, ctx)
	assert.NoError(t, err)
	fVersion = false
}

func TestVersion(t *testing.T) {
	Version()
}

func TestSelectInput_Bad(t *testing.T) {
	_, err := SelectInput("-")
	assert.Error(t, err)
}

func TestSelectInput_Good(t *testing.T) {
	fType = ".xml"
	r, err := SelectInput("-")
	assert.NoError(t, err)
	assert.EqualValues(t, os.Stdin, r)
	fType = ""
}

func TestSelectInput_Badfn(t *testing.T) {
	_, err := SelectInput("/testdata/google.com!keltia.net!1538438400!1538524798.xml")
	assert.Error(t, err)
}

// realmain() is the thing now
func TestMain_Noargs(t *testing.T) {
	assert.Error(t, realmain([]string{}))
}

func TestMain_Noargs_Verbose(t *testing.T) {
	fVerbose = true
	assert.Error(t, realmain([]string{}))
	fVerbose = false
}

func TestMain_Noargs_Debug(t *testing.T) {
	fDebug = true
	assert.Error(t, realmain([]string{}))
	fDebug = false
}

func TestMain_Noargs_NoResolv(t *testing.T) {
	fNoResolv = true
	r := realmain([]string{"testdata/google.com!keltia.net!1538438400!1538524799.xml"})
	fNoResolv = false
	assert.Empty(t, r)
}

func TestRealmain_GoodFile(t *testing.T) {
	r := realmain([]string{"testdata/google.com!keltia.net!1538438400!1538524799.xml"})
	assert.Empty(t, r)
}

func TestRealmain_GoodFile1(t *testing.T) {
	r := realmain([]string{"testdata/google.com!keltia.net!1538438400!1538524799.zip"})
	assert.Empty(t, r)
}

func TestRealmain_NoFile(t *testing.T) {
	r := realmain([]string{"/nonexistent"})
	assert.NotEmpty(t, r)
}

func TestMain_EmptyArg(t *testing.T) {
	r := realmain([]string{"foo"})
	assert.NotEmpty(t, r)
}

func TestMain_GoodFile(t *testing.T) {
	os.Args = append(os.Args, "testdata/google.com!keltia.net!1538438400!1538524799.zip")
	main()
}
