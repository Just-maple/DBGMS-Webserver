package logger

import (
	"fmt"
	"github.com/op/go-logging"
	"os"
)

const (
	DefaultLogFormat = `%{color}%{time:2006-01-02 15:04:05} %{shortfile} %{longfunc} ▶ [%{level:.4s}:%{id:03x}%{color:reset}] %{message}`
)

var Log *logging.Logger

func init() {
	Log = GetLogger("Server-debug", "./server-debug.log", logging.DEBUG)
}

func GetLogger(module, logfile string, level logging.Level) (logger *logging.Logger) {
	logger = logging.MustGetLogger(module)
	var backends []logging.Backend

	format := logging.MustStringFormatter(DefaultLogFormat)
	backendStd := logging.NewLogBackend(os.Stdout, fmt.Sprintf("[%s] ", module), 0)
	backendFormatter := logging.NewBackendFormatter(backendStd, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(level, "")
	backends = append(backends, backendLeveled)

	if logfile != "" {
		fd, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Errorf("os.Open(%s) err(%v)", logfile, err)
		}
		backendFile := logging.NewLogBackend(fd, "", 0)
		backendFormatter2 := logging.NewBackendFormatter(backendFile, format)
		backendLeveled2 := logging.AddModuleLevel(backendFormatter2)
		backendLeveled.SetLevel(level, "")

		backends = append(backends, backendLeveled2)
	}

	backendMulti := logging.MultiLogger(backends...)
	logger.ExtraCalldepth += 1
	logger.SetBackend(backendMulti)
	return logger
}
