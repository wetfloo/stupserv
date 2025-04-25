package common

import (
	"fmt"
	"os"
)

type ExitCode = int

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
	ExitOk ExitCode = iota
	ExitNotADir
	ExitInvalidArgs
	ExitFileErr
	ExitHttpServerErr
)

func PrintAndExit(
	code ExitCode,
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
