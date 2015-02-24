package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	path := flag.String("dir", ".", "write XML files to directory")
	flag.Parse()

	var err error
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()
	tee := io.TeeReader(os.Stdin, os.Stdout)
	report, err := Parse(tee)
	if err != nil {
		return
	}

	// Create a JUnit test suite xml report for each package.
	for _, pkg := range report.Packages {
		var f *os.File
		if f, err = os.Create(filepath.Join(*path, shortName(pkg.Name)+".xml")); err != nil {
			return
		}
		defer func() { err = f.Close() }()
		if err = JUnitReportXML(pkg, f); err != nil {
			return
		}
	}
}
