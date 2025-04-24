package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type argumentKey int
type argumentMap = map[string]argumentKey
type exitCode = int

type argValues struct {
	noCache bool
	addr    string
	path    string
}

const usageMsg = `
usage:
	stupserv [OPTIONS] [PATH]

PATH:
	used to determine what directory will be served. Default is the current working directory

OPTIONS:
	-a, --addr:     address and port to listen on. By default set to :6040
	-n, --no-cache: disables http cache. By default cache is used
	-h, --help:     displays this message
`

const (
	argAddr argumentKey = iota
	argNoCache
	argHelp
)

const (
	exitOk exitCode = iota
	exitNotADir
	exitInvalidArgs
	exitFileErr
	exitHttpServerErr
)

var arguments argumentMap = argumentMap{
	"-a":         argAddr,
	"--addr":     argAddr,
	"-n":         argNoCache,
	"--no-cache": argNoCache,
	"-h":         argHelp,
	"--help":     argHelp,
}

func main() {
	args := parseArgs()
	fileInfo, err := os.Stat(args.path)
	if err != nil {
		printAndExit(exitFileErr, false, "%v\n", err)
	}
	if !fileInfo.IsDir() {
		printAndExit(exitNotADir, true, "%v\n", err)
	}

	fileServer := http.FileServer(http.Dir(args.path))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if args.noCache {
			w.Header().Set("Cache-Control", "no-cache")
		}
		fileServer.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(args.addr, nil); err != nil {
		printAndExit(exitHttpServerErr, false, "%v\n", err)
	}
}

func parseArgs() argValues {
	args := os.Args[1:]
	var result argValues

	if len(args) >= 1 {
		result.path = args[len(args)-1]
	} else {
		ex, err := os.Executable()
		if err != nil {
			printAndExit(exitFileErr, false, "%v\n", err)
		}
		result.path = filepath.Dir(ex)
	}

	for i, arg := range args[:len(args)-1] {
		argKey, ok := arguments[arg]
		if !ok {
			continue
		}

		switch argKey {
		case argAddr:
			if i < len(args)-1 {
				result.addr = args[i+1]
				continue
			} else {
				printAndExit(exitInvalidArgs, true, "malformed args: couldn't parse address to listen on\n")
			}
		case argNoCache:
			result.noCache = true
		case argHelp:
			printAndExit(exitOk, true, "")
		}
	}

	if len(result.addr) == 0 {
		result.addr = ":6040"
	}

	return result
}

func printAndExit(
	code exitCode,
	showUsage bool,
	s string,
	f ...any,
) {
	fmt.Fprintf(os.Stderr, s, f...)
	if showUsage {
		fmt.Fprint(os.Stderr, usageMsg)
	}
	os.Exit(code)
}
