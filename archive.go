package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Extracter interface {
	Extract(t string) ([]byte, error)
	Close() error
}

type Plain struct {
	Name string
}

func (a Plain) Extract(t string) ([]byte, error) {
	ext := filepath.Ext(a.Name)
	if ext == t || t == "" {
		return ioutil.ReadFile(a.Name)
	}
	return []byte{}, fmt.Errorf("wrong file type")
}

func (a Plain) Close() error {
	return nil
}

type Zip struct {
	fn  string
	zfh *zip.ReadCloser
}

func NewZipfile(fn string) (*Zip, error) {
	zfh, err := zip.OpenReader(fn)
	if err != nil {
		return &Zip{}, errors.Wrap(err, "archive/zip")
	}
	return &Zip{fn: fn, zfh: zfh}, nil
}

func (a Zip) Extract(t string) ([]byte, error) {
	verbose("exploring %s", a.fn)

	for _, fn := range a.zfh.File {
		verbose("looking at %s", fn.Name)

		if path.Ext(fn.Name) == t ||
			path.Ext(fn.Name) == strings.ToUpper(t) {
			file, err := fn.Open()
			if err != nil {
				return []byte{}, errors.Wrapf(err, "no file matching type %s", t)
			}
			return ioutil.ReadAll(file)
		}
	}

	return []byte{}, fmt.Errorf("no file matching type %s", t)
}

func (a Zip) Close() error {
	return a.zfh.Close()
}

type Gzip struct {
	fn  string
	unc string
}

func NewGzipfile(fn string) (*Gzip, error) {
	base := filepath.Base(fn)
	pc := strings.Split(base, ".")
	unc := strings.Join(pc[0:len(pc)-1], ".")

	return &Gzip{fn: fn, unc: unc}, nil
}

func (a Gzip) Extract(t string) ([]byte, error) {
	buf, err := ioutil.ReadFile(a.fn)
	if err != nil {
		return []byte{}, errors.Wrap(err, "gzip/extract")
	}
	bufr := bytes.NewBuffer(buf)
	zfh, err := gzip.NewReader(bufr)
	defer zfh.Close()

	return ioutil.ReadAll(zfh)
}

func (a Gzip) Close() error {
	return nil
}

type Tar struct {
	fn string
}

func (a Tar) Extract(t string) ([]byte, error) {
	return []byte{}, nil
}

func (a Tar) Close() error {
	return nil
}

func NewArchive(fn string) (Extracter, error) {
	if fn == "" {
		return &Plain{}, fmt.Errorf("null string")
	}
	ext := filepath.Ext(fn)
	switch ext {
	case ".zip":
		return NewZipfile(fn)
	case ".gz":
		return NewGzipfile(fn)
	case ".tar":
		return &Tar{fn: fn}, nil
	}
	return &Plain{fn}, nil
}
