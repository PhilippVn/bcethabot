package internalerror

// Internal Error Module
// Used to print detailed internal Errors that shall not be sent to the user.
// source: https://snippets.aktagon.com/snippets/795-logging-in-golang-including-line-numbers-
import (
	"log"
	"runtime"
	"strings"
)

// Set this to true to enable DebugMode
var DebugMode bool = false

// Info example:
//
// Info("timezone %s", timezone)
//
func Info(msg string, vars ...interface{}) {
	log.Printf(strings.Join([]string{"[INFO ]", msg}, " "), vars...)
}

// Debug example:
//
// Debug("timezone %s", timezone)
//
func Debug(msg string, vars ...interface{}) {
	if DebugMode {
		log.Printf(strings.Join([]string{"[DEBUG]", msg}, " "), vars...)
	}
}

// Fatal example:
//
// Fatal(errors.New("db timezone must be UTC"))
//
func Fatal(err error) {
	pc, fn, line, _ := runtime.Caller(1)
	// Include function name if debugging
	if DebugMode {
		log.Fatalf("[FATAL] %s [%s:%s:%d]", err, runtime.FuncForPC(pc).Name(), fn, line)
	} else {
		log.Fatalf("[FATAL] %s [%s:%d]", err, fn, line)
	}
}

// Error example:
//
// Error(errors.Errorf("Invalid timezone %s", timezone))
//
func Error(err error) {
	pc, fn, line, _ := runtime.Caller(1)
	// Include function name if debugging
	if DebugMode {
		log.Printf("[ERROR] %s [%s:%s:%d]", err, runtime.FuncForPC(pc).Name(), fn, line)
	} else {
		log.Printf("[ERROR] %s [%s:%d]", err, fn, line)
	}
}
