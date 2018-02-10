package log

import (
	"os"

	"github.com/crazy-max/nodejs-portable/app/fs"
	"github.com/crazy-max/nodejs-portable/app/pathu"
	logging "github.com/op/go-logging"
)

// Logger
var (
	Logger    = logging.MustGetLogger("nodejs-portable")
	Logfile   *os.File
	logFormat = logging.MustStringFormatter(`%{time:2006-01-02 15:04:05} %{level:.4s} - %{message}`)
)

func init() {
	var err error

	// log file
	pathu.CurrentPath = fs.FormatWinPath(pathu.CurrentPath)
	Logfile, err = fs.OpenFile(fs.Join(pathu.CurrentPath, "nodejs-portable.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		Logger.Error("Log file:", err)
	}

	// init logger
	logBackendFile := logging.NewBackendFormatter(logging.NewLogBackend(Logfile, "", 0), logFormat)
	logging.SetBackend(logBackendFile)
}
