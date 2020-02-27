package common

import (
	"runtime"
	"strings"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/cosmos/cosmos-sdk/client/context"
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

func RobustQuery(cliCtx context.CLIContext, route string) (res []byte, err error) {
	res, _, err = cliCtx.Query(route)
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		res, _, err = cliCtx.Query(route)
		return
	}

	return
}

func RobustQueryWithData(cliCtx context.CLIContext, route string, data []byte) (res []byte, err error) {
	res, _, err = cliCtx.QueryWithData(route, data)
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		res, _, err = cliCtx.QueryWithData(route, data)
		return
	}

	return
}
