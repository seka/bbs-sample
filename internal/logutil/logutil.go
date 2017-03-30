package logutil

import (
	"os"

	"github.com/inconshreveable/log15"
)

// SeetupRootLogger ...
func SeetupRootLogger(logLevel string) {
	lvl, err := log15.LvlFromString(logLevel)
	if err != nil {
		log15.Error("LvlFromString error", "err", err)
		os.Exit(1)
	}
	h := log15.LvlFilterHandler(lvl, log15.StdoutHandler)
	log15.Root().SetHandler(h)
}
