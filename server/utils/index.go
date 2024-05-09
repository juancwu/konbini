package utils

import (
	"os"

	"github.com/charmbracelet/log"
)

var logger *log.Logger

func init() {
	logger = log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		TimeFormat:      "2006/01/02 15:04:05.0.000000000",
	})
}

func Logger() *log.Logger {
	return logger
}