package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bruth/flatjson"
)

var usage = `usage: flatjson [options] [path]

flatjson takes a JSON string and re-encodes into a flat map or array of
key-value pairs.

Examples:

  Read from stdin:

    cat file.json | flatjson

  Read from a file, output as an array instead of a map:

    flatjson -array file.json

Options:

`

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
		flag.PrintDefaults()
	}
}

func main() {
	var array bool

	flag.BoolVar(&array, "array", false, "Output as an array of pairs.")
	flag.Parse()

	args := flag.Args()

	var r io.Reader

	// Use stdin if no path is supplied.
	if len(args) == 0 {
		r = os.Stdin
	} else {
		f, err := os.Open(args[0])

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		r = f
	}

	enc := flatjson.NewEncoder(os.Stdout)

	var err error

	if array {
		err = enc.ConvertArray(r)
	} else {
		err = enc.ConvertMap(r)
	}

	if err != nil {
		log.Fatal(err)
	}
}
