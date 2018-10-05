package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/keltia/sandbox"
)

var (
	// MyName is the application
	MyName = filepath.Base(os.Args[0])
	// MyVersion is our version
	MyVersion = "0.4.0-wip"
	// Author should be abvious
	Author = "Ollivier Robert"

	fDebug   bool
	fVerbose bool
)

func init() {
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
}

func main() {
	flag.Parse()

	if fDebug {
		fVerbose = true
		debug("debug mode")
	}

	if len(flag.Args()) != 1 {
		log.Fatalln("You must specify at least one file.")
	}

	snd, err := sandbox.New(MyName)
	if err != nil {
		log.Fatalf("Fatal: Can not create sandbox: %v", err)
	}
	defer snd.Cleanup()

	file := flag.Arg(0)
	txt, err := HandleSingleFile(snd, file)
	if err != nil {
		log.Fatalf("error handling %s: %v", file, err)
	}
	fmt.Println(txt)
}
