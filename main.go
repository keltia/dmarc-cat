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
	MyVersion = "0.3.0"
	// Author should be abvious
	Author = "Ollivier Robert"

	fDebug   bool
	fVerbose bool

	tempdir string
)

func init() {
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
}

func main() {
	flag.Parse()

	if fDebug {
		fVerbose = true
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
	err = snd.Run(func() error {
		var err error

		if text, err := HandleSingleFile(snd.Cwd(), file); err != nil {
			log.Printf("error parsing %s: %v", file, err)
		} else {
			fmt.Println(text)
		}
		return err
	})
	if err != nil {
		log.Printf("error handling %s: %v", file, err)
	}
}
