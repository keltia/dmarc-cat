package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	// MyName is the application
	MyName = filepath.Base(os.Args[0])
	// MyVersion is our version
	MyVersion = "0.7.0"
	// Author should be abvious
	Author = "Ollivier Robert"

	fDebug    bool
	fNoResolv bool
	fVerbose  bool
)

// Context is passed around rather than being a global var/struct
type Context struct {
	r Resolver
}

func init() {
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fNoResolv, "N", false, "Do not resolve IPs")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
}

func main() {
	flag.Parse()

	if fDebug {
		fVerbose = true
		debug("debug mode")
	}

	if len(flag.Args()) != 1 {
		log.Println("You must specify at least one file.")
		return
	}

	ctx := &Context{RealResolver{}}

	// Make it easier to sub it out
	if fNoResolv {
		ctx.r = NullResolver{}
	}

	file := flag.Arg(0)
	txt, err := HandleSingleFile(ctx, file)
	if err != nil {
		log.Printf("error handling %s: %v", file, err)
		return
	}
	fmt.Println(txt)
}
