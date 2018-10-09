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

// Extracter is the main interface we have
type Extracter interface {
	Extract(t string) ([]byte, error)
}

// ExtractCloser is the same with Close()
type ExtractCloser interface {
	Extracter
	Close() error
}

// Plain is for plain text
type Plain struct {
	Name string
}

// Extract returns the content of the file
func (a Plain) Extract(t string) ([]byte, error) {
	ext := filepath.Ext(a.Name)
	if ext == t || t == "" {
		return ioutil.ReadFile(a.Name)
	}
	return []byte{}, fmt.Errorf("wrong file type")
}

// Close is a no-op
func (a Plain) Close() error {
	return nil
}

// Zip is for pkzip/infozip files
type Zip struct {
	fn  string
	zfh *zip.ReadCloser
}

// NewZipfile open the zip file
func NewZipfile(fn string) (*Zip, error) {
	zfh, err := zip.OpenReader(fn)
	if err != nil {
		return &Zip{}, errors.Wrap(err, "archive/zip")
	}
	return &Zip{fn: fn, zfh: zfh}, nil
}

// Extract returns the content of the file
func (a Zip) Extract(t string) ([]byte, error) {
	verbose("exploring %s", a.fn)

	ft := strings.ToLower(t)
	for _, fn := range a.zfh.File {
		verbose("looking at %s", fn.Name)

		if path.Ext(fn.Name) == ft {
			file, err := fn.Open()
			if err != nil {
				return []byte{}, errors.Wrapf(err, "no file matching type %s", t)
			}
			return ioutil.ReadAll(file)
		}
	}

	return []byte{}, fmt.Errorf("no file matching type %s", t)
}

// Close does something here
func (a Zip) Close() error {
	return a.zfh.Close()
}

// Gzip is a gzip-compressed file
type Gzip struct {
	fn  string
	unc string
}

// NewGzipfile stores the uncompressed file name
func NewGzipfile(fn string) (*Gzip, error) {
	base := filepath.Base(fn)
	pc := strings.Split(base, ".")
	unc := strings.Join(pc[0:len(pc)-1], ".")

	return &Gzip{fn: fn, unc: unc}, nil
}

// Extract returns the content of the file
func (a Gzip) Extract(t string) ([]byte, error) {
	buf, err := ioutil.ReadFile(a.fn)
	if err != nil {
		return []byte{}, errors.Wrap(err, "gzip/extract")
	}
	bufr := bytes.NewBuffer(buf)
	zfh, err := gzip.NewReader(bufr)
	if err != nil {
		return []byte{}, errors.Wrap(err, "gunzip")
	}
	content, err := ioutil.ReadAll(zfh)
	defer zfh.Close()

	return content, err
}

// Close is a no-op
func (a Gzip) Close() error {
	return nil
}

// NewArchive is the main creator
func NewArchive(fn string) (ExtractCloser, error) {
	if fn == "" {
		return &Plain{}, fmt.Errorf("null string")
	}
	ext := filepath.Ext(fn)
	switch ext {
	case ".zip":
		return NewZipfile(fn)
	case ".gz":
		return NewGzipfile(fn)
	}
	return &Plain{fn}, nil
}
