go-junit-report
===============

Converts `go test -v` output to xml reports, suitable for applications that
expect junit xml reports. (e.g. [Jenkins](http://jenkins-ci.org)), one per Go package.

Installation
------------

	go get github.com/cee-dub/go-junit-report

Usage
-----

	mkdir -p junit
	go test -v | go-junit-report -dir junit

