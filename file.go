package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

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

// OpenFile looks at the file and give it to openZipfile() if needed
func OpenFile(tempdir, file string) (r io.ReadCloser, err error) {
	var myfile string

	if _, err = os.Stat(file); err != nil {
		return nil, errors.Wrap(err, "OpenFile/stat")
	}

	myfile = file

	// Next pass, check for zip file
	if path.Ext(myfile) == ".zip" ||
		path.Ext(myfile) == ".ZIP" {

		verbose("found zip file %s", myfile)

		myfile = openZipfile(tempdir, myfile)
	} else if path.Ext(myfile) == ".gz" ||
		path.Ext(myfile) == ".GZ" {
		verbose("found gzip file %s", myfile)

		myfile = OpenGzipFile(tempdir, myfile)
	}
	return os.Open(myfile)
}

// ExtractXML reads the first xml in the zip file and copy into a temp file
func ExtractXML(tempdir string, fn *zip.File) (file string) {
	verbose("found %s", fn.Name)

	// Open the stream
	fh, err := fn.Open()
	if err != nil {
		log.Fatalf("unable to extract %s", fn.Name)
	}

	// Create our temp file
	ours, err := os.Create(filepath.Join(tempdir, fn.Name))
	if err != nil {
		log.Fatalf("unable to create %s in %s: %v", fn.Name, tempdir, err)
	}
	defer ours.Close()

	verbose("created our tempfile %s", filepath.Join(tempdir, fn.Name))

	// copy all the bits over
	_, err = io.Copy(ours, fh)
	if err != nil {
		log.Fatalf("unable to write %s in %s: %v", fn.Name, tempdir, err)
	}
	file = filepath.Join(tempdir, fn.Name)
	return
}

// openGzipfile uncompress the file and store it into a .csv file in sandbox
func OpenGzipFile(tempdir, file string) (fname string) {

	buf, err := ioutil.ReadFile(file)
	bufr := bytes.NewBuffer(buf)
	zfh, err := gzip.NewReader(bufr)
	if err != nil {
		log.Fatalf("error opening %s: %v", file, err)
	}
	defer zfh.Close()

	verbose("exploring %s", file)

	file = filepath.Base(file)

	cmps := strings.Split(file, ".")
	if cmps == nil {
		log.Fatalf("error, file not csv: %s", file)
	}

	file = strings.Join(cmps[0:len(cmps)-1], ".")
	// Create our temp file
	ours, err := os.Create(filepath.Join(tempdir, file))
	if err != nil {
		log.Fatalf("unable to create %s in %s: %v", file, tempdir, err)
	}
	defer ours.Close()

	verbose("created our tempfile %s", filepath.Join(tempdir, file))

	// copy all the bits over
	_, err = io.Copy(ours, zfh)

	fname = file
	return
}

// openZipfile extracts the first XML file out of he given zip.
func openZipfile(tempdir, file string) (fname string) {
	// Go on
	if err := os.Chdir(tempdir); err != nil {
		log.Fatalf("unable to use tempdir %s: %v", tempdir, err)
	}

	zfh, err := zip.OpenReader(file)
	if err != nil {
		log.Fatalf("error opening %s: %v", file, err)
	}
	defer zfh.Close()

	verbose("exploring %s", file)

	for _, fn := range zfh.File {
		verbose("looking at %s", fn.Name)

		if path.Ext(fn.Name) == ".xml" ||
			path.Ext(fn.Name) == ".XML" {
			file = ExtractXML(tempdir, fn)
			break
		}
	}
	fname = file
	return
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

	var fh io.ReadCloser

	err = snd.Run(func() error {
		var err error

		tempdir := snd.Cwd()
		// Look at the file and whatever might be inside (and decrypt/unzip/…)
		fh, err = OpenFile(tempdir, myfile)
		if err != nil {
			return errors.Wrap(err, "OpenFile")
		}
		return err
	})

	debug("fh=%#v", fh)
	verbose("Analyzing %s", myfile)
	report, err := ParseXML(fh)
	if err != nil {
		return "", errors.Wrap(err, "error parsing XML")
	}

	debug("report=%v\n", report)

	output, err := Analyze(report)
	if err != nil {
		log.Printf("error: %v", err)
	}

	return output, nil
}
