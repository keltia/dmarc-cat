package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var (
	// MyName is the application
	MyName = "dmarc-cat"
	// MyVersion is our version
	MyVersion = "0.1.0"

	fDebug   bool
	fVerbose bool

	tempdir string
)

func init() {
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
}

// cleanupTemp removes the temporary directory
func cleanupTemp(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		log.Printf("cleanup failed for %s: %v", dir, err)
	}
}

// createSandbox creates our our directory with TEMPDIR (wherever it is)
func createSandbox(tag string) (path string) {
	// Extract in safe location
	dir, err := ioutil.TempDir("", tag)
	if err != nil {
		log.Fatalf("unable to create sandbox %s: %v", dir, err)
	}
	return dir
}

func main() {
	flag.Parse()

	if fDebug {
		fVerbose = true
	}

	if len(flag.Args()) != 1 {
		log.Fatalln("You must specify at least one file.")
	}

	tempdir = createSandbox(MyName)
	if err := handleSingleFile(tempdir, flag.Arg(0)); err != nil {
		log.Printf("error parsing %s: %v", flag.Arg(0), err)
	}
	cleanupTemp(tempdir)
}
