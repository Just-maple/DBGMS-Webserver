package logger

import (
	"github.com/op/go-logging"
	"logs"
)

var Log *logging.Logger

func init() {
	Log = logs.GetLogger("WebServer", "./server-debug.log", logging.DEBUG)
}
