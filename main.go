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
	cache bool
	addr  string
	path  string
}

const usageMsg = `
usage:
	stupserv [OPTIONS] [PATH]

PATH:
	used to determine what directory will be served. Default is the current working directory

OPTIONS:
	-a, --addr:     address and port to listen on. By default set to :6040
	-c, --cache:    enables http cache. By default cache is not used
	-h, --help:     displays this message
`

const (
	argAddr argumentKey = iota
	argCache
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
	"-a":      argAddr,
	"--addr":  argAddr,
	"-c":      argCache,
	"--cache": argCache,
	"-h":      argHelp,
	"--help":  argHelp,
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
		if !args.cache {
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

	for i, arg := range args[:max(len(args)-1, 0)] {
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
		case argCache:
			result.cache = true
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
