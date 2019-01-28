package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/keltia/archive"
)

var (
	// MyName is the application
	MyName = filepath.Base(os.Args[0])
	// MyVersion is our version
	MyVersion = "0.9.2"
	// Author should be abvious
	Author = "Ollivier Robert"

	fDebug    bool
	fJobs     int
	fNoResolv bool
	fSort     string
	fVerbose  bool
	fVersion  bool
)

// Context is passed around rather than being a global var/struct
type Context struct {
	r    Resolver
	jobs int
}

func init() {
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fNoResolv, "N", false, "Do not resolve IPs")
	flag.IntVar(&fJobs, "j", runtime.NumCPU(), "Parallel jobs")
	flag.StringVar(&fSort, "S", `"Count" "dsc"`, "Sort results")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
	flag.BoolVar(&fVersion, "version", false, "Display version")
}

func Version() {
	fmt.Printf("%s version %s/j%d archive/%s\n", MyName, MyVersion, fJobs, archive.Version())
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

	ctx := &Context{RealResolver{}, fJobs}

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
