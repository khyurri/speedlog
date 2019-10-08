package utils

import (
	"log"
	"path/filepath"
	"runtime"
)

const (
	LG_ERROR = iota // show only errors (default)
	LG_DEBUG        // show debug messages
)

var Level = LG_ERROR

// Ok — checks err and prints message if err is not nil
func Ok(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
	}
}

// Fatal is equivalent Ok followed by a call to os.Exit(1).
func Fatal(err error) {

}

// Panic is equivalent Ok followed by a call to panic().
func Panic(err error) {

}

// Debug prints debug message if flag --debug have been passed
func Debug(msg interface{}) {
	if Level == LG_DEBUG {
		log.Printf("[DEBUG] %+v\n", msg)
	}
}
