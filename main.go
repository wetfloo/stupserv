package main

import (
	"net/http"
	"os"

	"github.com/wetfloo/stupserv/args"
	"github.com/wetfloo/stupserv/common"
)

func main() {
	args := args.ParseArgs()
	fileInfo, err := os.Stat(args.Path)
	if err != nil {
		common.PrintAndExit(common.ExitFileErr, false, "%v\n", err)
	}
	if !fileInfo.IsDir() {
		common.PrintAndExit(common.ExitNotADir, true, "%v\n", err)
	}

	fileServer := http.FileServer(http.Dir(args.Path))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !args.Cache {
			w.Header().Set("Cache-Control", "no-cache")
		}
		fileServer.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(args.Addr, nil); err != nil {
		common.PrintAndExit(common.ExitHttpServerErr, false, "%v\n", err)
	}
}
