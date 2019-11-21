package common

import (
	"runtime"
	"strings"

	"github.com/celer-network/goutils/log"
)

// EnableLogLongFile set the log file splitter from the sgn root folder
func EnableLogLongFile() {
	log.EnableLongFile()
	_, file, _, ok := runtime.Caller(0)
	if ok {
		pref := file[:strings.LastIndex(file[:strings.LastIndex(file, "/")], "/")+1]
		log.SetFilePathSplit(pref)
	}
}
