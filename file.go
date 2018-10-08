package main

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/keltia/sandbox"
	"github.com/pkg/errors"
)

const (
	reFileName = `^([\w\.]+)!([\w\.]+)!([\d]+)!([\d]+)\.(xml\.gz|zip)$`
)

func checkFilename(file string) (ok bool) {
	base := filepath.Base(file)
	re := regexp.MustCompile(reFileName)

	return re.MatchString(base)
}

// HandleSingleFile creates a tempdir and dispatch csv/zip files to handler.
func HandleSingleFile(snd *sandbox.Dir, file string) (string, error) {
	var myfile string

	debug("file=%s", file)
	if !checkFilename(file) {
		return "", fmt.Errorf("bad filename")
	}

	// We want the full path
	myfile, err := filepath.Abs(file)
	if err != nil {
		return "", errors.Wrapf(err, "Abs(%s)", file)
	}

	a, err := NewArchive(myfile)
	body, err := a.Extract(".xml")

	debug("xml=%#v", body)
	verbose("Analyzing %s", myfile)

	var report Feedback

	if err := xml.Unmarshal(body, &report); err != nil {
		return "", errors.Wrap(err, "unmarshall")
	}

	debug("report=%v\n", report)

	return Analyze(report)
}
