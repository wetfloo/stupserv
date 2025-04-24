package main

import (
	"fmt"
	"os"
)

type argumentKey int
type argumentMap = map[string]argumentKey

type argValues struct {
	addr string
	path string
}

const usageMsg = `
usage:
	stupserv [OPTIONS] [PATH]

PATH:
	used to determine what directory will be served. Default is the current working directory

OPTIONS:
	-a, --addr: address and port to listen on. By default set to :6040
	-h, --help: displays this message
`

const (
	argAddr argumentKey = iota
	argHelp
)

const (
	exitOk int = iota
	exitNotADir
	exitInvalidArgs
	exitFileErr
)

var arguments argumentMap = argumentMap{
	"-a":     argAddr,
	"--addr": argAddr,
	"-h":     argHelp,
	"--help": argHelp,
}

func main() {
	args := parseArgs()
	fileInfo, err := os.Stat(args.path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(exitFileErr)
	}
	if !fileInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "provided path %s is not a directory\n")
		os.Exit(exitNotADir)
	}
}

func parseArgs() argValues {
	args := os.Args[1:]
	var result argValues

	if len(args) > 1 {
		result.path = os.Args[len(args)-1]
	} else {
		p, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(exitFileErr)
		}
		result.path = p
	}

	for i, arg := range args {
		argKey, ok := arguments[arg]
		if !ok {
			if arg == result.path {
				continue
			}
			if _, err := os.Stat(arg); err != nil {
				fmt.Fprintf(os.Stderr, "path exists for the argument %s, but path had already been set, and is %s\n", result.path)
				fmt.Fprint(os.Stderr, usageMsg)
				os.Exit(exitInvalidArgs)
			}

			fmt.Fprintf(os.Stderr, "arg %s is not a valid argument\n", argKey)
			fmt.Fprint(os.Stderr, usageMsg)
			os.Exit(exitInvalidArgs)
		}

		switch argKey {
		case argAddr:
			if i < len(args)-1 {
				result.addr = args[i+1]
				continue
			} else {
				fmt.Fprint(os.Stderr, "malformed args: couldn't parse address to listen on\n")
				fmt.Fprint(os.Stderr, usageMsg)
				os.Exit(exitInvalidArgs)
			}
		case argHelp:
			fmt.Fprint(os.Stderr, usageMsg)
			os.Exit(0)
		}
	}

	return result
}
