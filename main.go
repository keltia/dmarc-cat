package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/keltia/archive"
)

var (
	// MyName is the application
	MyName = filepath.Base(os.Args[0])
	// MyVersion is our version
	MyVersion = "0.8.0"
	// Author should be abvious
	Author = "Ollivier Robert"

	fDebug    bool
	fNoResolv bool
	fVerbose  bool
	fVersion  bool
)

// Context is passed around rather than being a global var/struct
type Context struct {
	r Resolver
}

func init() {
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fNoResolv, "N", false, "Do not resolve IPs")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
	flag.BoolVar(&fVersion, "version", false, "Display version")
}

func Version() {
	fmt.Printf("%s version %s archive/%s\n", MyName, MyVersion, archive.Version())
}

// Setup creates our context and check stuff
func Setup(a []string) *Context {
	// Exist early if -version
	if fVersion {
		Version()
		return nil
	}

	if fDebug {
		fVerbose = true
		debug("debug mode")
	}

	if len(a) < 1 {
		log.Println("You must specify at least one file.")
		return nil
	}

	ctx := &Context{RealResolver{}}

	// Make it easier to sub it out
	if fNoResolv {
		ctx.r = NullResolver{}
	}

	return ctx
}

func main() {
	flag.Parse()

	ctx := Setup(flag.Args())
	if ctx == nil {
		return
	}

	file := flag.Arg(0)
	txt, err := HandleSingleFile(ctx, file)
	if err != nil {
		log.Printf("error handling %s: %v", file, err)
		return
	}
	fmt.Println(txt)
}
