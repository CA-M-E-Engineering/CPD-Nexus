package logger

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger
)

func init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Infof(format string, v ...interface{}) {
	Info.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	Error.Printf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	Debug.Printf(format, v...)
}
