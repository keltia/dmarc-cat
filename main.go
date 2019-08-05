package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/keltia/archive"
	"github.com/pkg/errors"
)

var (
	// MyName is the application
	MyName = filepath.Base(os.Args[0])
	// MyVersion is our version
	MyVersion = "0.12.0,parallel"
	// Author should be obvious
	Author = "Ollivier Robert"

	fDebug    bool
	fJobs     int
	fNoResolv bool
	fSort     string
	fType     string
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
	flag.StringVar(&fType, "t", "", "File type for stdin mode")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
	flag.BoolVar(&fVersion, "version", false, "Display version")
}

func Version() {
	fmt.Printf("%s version %s/j%d archive/%s\n", MyName, MyVersion, fJobs, archive.Version())
}

// Setup creates our context and check stuff
func Setup(a []string) (*Context, error) {
	// Exist early if -version
	if fVersion {
		Version()
		return nil, nil
	}

	if fDebug {
		fVerbose = true
		debug("debug mode")
	}

	if len(a) < 1 {
		return nil, fmt.Errorf("You must specify at least one file.")
	}

	ctx := &Context{RealResolver{}, fJobs}

	// Make it easier to sub it out
	if fNoResolv {
		ctx.r = NullResolver{}
	}

	return ctx, nil
}

func SelectInput(file string) (io.ReadCloser, error) {
	debug("file=%s", file)
	debug("file=%s", file)

	if file == "-" {
		if fType == "" {
			return nil, errors.New("Wrong file type, use -t")
		}
		return os.Stdin, nil
	}

	// We have a filename
	if !checkFilename(file) {
		return nil, errors.New("bad filename")
	}

	// We want the full path
	myfile, err := filepath.Abs(file)
	if err != nil {
		return nil, errors.Wrapf(err, "Abs(%s)", file)
	}

	return os.Open(myfile)
}

func realmain(args []string) error {
	flag.Parse()

	ctx, err := Setup(args)
	if ctx == nil {
		return errors.Wrap(err, "realmain")
	}

	var txt string

	// Look for input file or stdin/"-"
	file := args[0]

	verbose("Analyzing %s", file)

	if filepath.Ext(file) == ".zip" ||
		filepath.Ext(file) == ".ZIP" {

		txt, err = HandleZipFile(ctx, file)
		if err != nil {
			return errors.Wrapf(err, "file %s:", file)
		}
	} else {
		in, err := SelectInput(file)
		if err != nil {
			return errors.Wrap(err, "SelectInput")
		}
		defer in.Close()

		typ := archive.Ext2Type(fType)

		txt, err = HandleSingleFile(ctx, in, typ)
		if err != nil {
			return errors.Wrapf(err, "file %s:", file)
		}
	}
	fmt.Println(txt)
	return nil
}

func main() {
	// Parse CLI
	flag.Parse()

	if err := realmain(flag.Args()); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
