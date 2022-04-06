package main

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/keltia/archive"
	"github.com/pkg/errors"
)

/* cf. https://tools.ietf.org/html/rfc7489#section-7.2.1.1

filename = receiver "!" policy-domain "!" begin-timestamp
            "!" end-timestamp [ "!" unique-id ] "." extension

unique-id = 1*(ALPHA / DIGIT)
*/
const (
	reFileName = `^([\S\.]+)!([\S\.]+)!([\d]+)!([\d]+)(![[:alnum:]]+)*(\.\S+)(\.(gz|zip))*$`
)

var reFN *regexp.Regexp

func init() {
	reFN = regexp.MustCompile(reFileName)
}

func checkFilename(file string) (ok bool) {
	base := filepath.Base(file)
	return reFN.MatchString(base)
}

// HandleZipFile is here for zip files because archive.NewFromReader() does not work here
func HandleZipFile(ctx *Context, file string) (string, error) {
	debug("HandleZipFile")

	var body []byte

	a, err := archive.New(file)
	if err == nil {

		body, err = a.Extract(".xml")
		if err != nil {
			return "", errors.Wrap(err, "extract")
		}
	} else {
		// Got plain text (i.e. xml)
		if body, err = ioutil.ReadFile(file); err != nil {
			return "", errors.Wrap(err, "ReadFile")
		}
	}

	debug("xml=%s", string(body))

	var report Feedback

	if err := xml.Unmarshal(body, &report); err != nil {
		return "", errors.Wrap(err, "unmarshall")
	}

	debug("report=%v\n", report)

	return Analyze(ctx, report)
}

// HandleSingleFile creates a tempdir and dispatch csv/zip files to handler.
func HandleSingleFile(ctx *Context, r io.ReadCloser, typ int) (string, error) {
	debug("HandleSingleFile")

	var body []byte

	debug("typ=%d", typ)
	if typ == archive.ArchiveZip {
		return "", errors.New("unsupported")
	}

	a, err := archive.NewFromReader(r, typ)
	if err == nil {

		debug("a=%#v", a)
		body, err = a.Extract("")
		if err != nil {
			return "", errors.Wrap(err, "extract")
		}
	} else {
		// Got plain text (i.e. xml)
		buf := bytes.NewBuffer(body)
		_, err := io.Copy(buf, r)
		if err != nil {
			return "", errors.Wrap(err, "copy")
		}
	}
	debug("xml=%#v", body)

	var report Feedback

	if err := xml.Unmarshal(body, &report); err != nil {
		debug("%d %s", typ, fType)
		return "", errors.Wrap(err, "unmarshall")
	}

	debug("report=%v\n", report)

	return Analyze(ctx, report)
}
