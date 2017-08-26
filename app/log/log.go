package log

import (
	"github.com/op/go-logging"
)

// Logger
var (
	log = logging.MustGetLogger("nodejs-portable")
)

// Error logs a message using ERROR as log level.
func Error(args ...interface{}) {
	log.Error(args)
}

// Info logs a message using INFO as log level.
func Info(args ...interface{}) {
	log.Info(args)
}
