package log

import (
	"fmt"
	"os"

	"github.com/tendermint/tendermint/libs/log"
)

var logger log.Logger

func init() {
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

// TODO: adjust to tendermint logger's key-value style
func Debug(v ...interface{}) {
	logger.Debug(fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	logger.Debug(fmt.Sprintf(format, v...))
}

func Debugln(v ...interface{}) {
	logger.Debug(fmt.Sprintln(v...))
}

func Info(v ...interface{}) {
	logger.Info(fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...))
}

func Infoln(v ...interface{}) {
	logger.Info(fmt.Sprintln(v...))
}

func Error(v ...interface{}) {
	logger.Error(fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...))
}

func Errorln(v ...interface{}) {
	logger.Error(fmt.Sprintln(v...))
}

// Fatal use "(Fatal)" to indicate this is a Fatal log
// since this logger only supports Debug/Info/Error
func Fatal(v ...interface{}) {
	logger.Error("(Fatal)", fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	logger.Error("(Fatal)", fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Fatalln(v ...interface{}) {
	logger.Error("(Fatal)", fmt.Sprintln(v...))
	os.Exit(1)
}
