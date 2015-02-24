package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"runtime"
	"strings"
)

// Types defining XML output compatible with the format described at
// http://windyroad.org/dl/Open%20Source/JUnit.xsd
type (
	JUnitTestSuite struct {
		XMLName    xml.Name        `xml:"testsuite"`
		Tests      int             `xml:"tests,attr"`
		Failures   int             `xml:"failures,attr"`
		Skips      int             `xml:"skips,attr"`
		Time       string          `xml:"time,attr"`
		Name       string          `xml:"name,attr"`
		Properties []JUnitProperty `xml:"properties>property,omitempty"`
		TestCases  []JUnitTestCase
	}

	JUnitTestCase struct {
		XMLName   xml.Name      `xml:"testcase"`
		Classname string        `xml:"classname,attr"`
		Name      string        `xml:"name,attr"`
		Time      string        `xml:"time,attr"`
		Failure   *JUnitFailure `xml:"failure,omitempty"`
		Skip      *JUnitSkip    `xml:"skipped,omitempty"`
	}

	JUnitProperty struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	}

	JUnitFailure struct {
		Message  string `xml:"message,attr"`
		Type     string `xml:"type,attr"`
		Contents string `xml:",chardata"`
	}

	JUnitSkip struct {
		Message  string `xml:"message,attr"`
		Type     string `xml:"type,attr"`
		Contents string `xml:",chardata"`
	}
)

func NewJUnitProperty(name, value string) JUnitProperty {
	return JUnitProperty{
		Name:  name,
		Value: value,
	}
}

// JUnitReportXML writes a junit xml representation of the given report to w, including an
// XML header.
func JUnitReportXML(pkg Package, w io.Writer) error {
	ts := JUnitTestSuite{
		Tests:      len(pkg.Tests),
		Time:       fmtDur(pkg.Time),
		Name:       pkg.Name,
		Properties: []JUnitProperty{NewJUnitProperty("go.version", runtime.Version())},
	}

	// individual test cases
	for _, test := range pkg.Tests {
		testCase := JUnitTestCase{
			Classname: shortName(pkg.Name),
			Name:      test.Name,
			Time:      fmtDur(test.Time),
		}
		switch test.Result {
		case FAIL:
			ts.Failures++
			testCase.Failure = &JUnitFailure{
				Message:  "Failed",
				Contents: strings.Join(test.Output, "\n"),
			}
		case SKIP:
			ts.Skips++
			testCase.Skip = &JUnitSkip{
				Message:  "Skipped",
				Contents: strings.Join(test.Output, "\n"),
			}

		}

		ts.TestCases = append(ts.TestCases, testCase)
	}

	enc := xml.NewEncoder(w)
	enc.Indent("", "\t")
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("writing xml header: %s", err)
	}
	if err := enc.Encode(ts); err != nil {
		return fmt.Errorf("writing xml: %s", err)
	}
	_, err := w.Write([]byte("\n"))
	return err
}

func fmtDur(time int) string {
	return fmt.Sprintf("%.3f", float64(time)/1000.0)
}

// shortname truncates name to the last element after a '/', if one occurs.
func shortName(name string) string {
	if idx := strings.LastIndex(name, "/"); idx > -1 && idx < len(name) {
		name = name[idx+1:]
	}
	return name
}
