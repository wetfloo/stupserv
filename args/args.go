package args

import (
	"os"
	"path/filepath"

	"github.com/wetfloo/stupserv/common"
)

type argumentKey int
type argumentMap = map[string]argumentKey

type Values struct {
	Cache bool
	Addr  string
	Path  string
}

const (
	argAddr argumentKey = iota
	argCache
	argHelp
)

var arguments argumentMap = argumentMap{
	"-a":      argAddr,
	"--addr":  argAddr,
	"-c":      argCache,
	"--cache": argCache,
	"-h":      argHelp,
	"--help":  argHelp,
}

func ParseArgs(args []string) Values {
	var result Values

	if len(args) >= 1 {
		result.Path = args[len(args)-1]
	} else {
		ex, err := os.Executable()
		if err != nil {
			common.PrintAndExit(common.ExitFileErr, false, "%v\n", err)
		}
		result.Path = filepath.Dir(ex)
	}

	for i, arg := range args[:max(len(args)-1, 0)] {
		argKey, ok := arguments[arg]
		if !ok {
			continue
		}

		switch argKey {
		case argAddr:
			if i < len(args)-1 {
				result.Addr = args[i+1]
				continue
			} else {
				common.PrintAndExit(
					common.ExitInvalidArgs,
					true,
					"malformed args: couldn't parse address to listen on\n",
				)
			}
		case argCache:
			result.Cache = true
		case argHelp:
			common.PrintAndExit(common.ExitOk, true, "")
		}
	}

	if len(result.Addr) == 0 {
		result.Addr = ":6040"
	}

	return result
}
