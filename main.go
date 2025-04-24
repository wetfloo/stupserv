package main

import (
	"fmt"
	"os"
)

const invalidArgsMsg = `
usage:
	stupserv [OPTIONS] [PATH]

PATH:
	used to determine what directory will be served. Default is the current working directory

OPTIONS:
	-a, --addr: address and port to listen on. By default set to :6040
	-h, --help: displays this message
`

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" {
			fmt.Fprint(os.Stderr, invalidArgsMsg)
			os.Exit(0)
		}
	}

	var path string
	if len(os.Args) > 2 {
		path = os.Args[len(os.Args)-1]
	} else {
		p, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}
		path = p
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	if !fileInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "provided path %s is not a directory")
		os.Exit(2)
	}
}
