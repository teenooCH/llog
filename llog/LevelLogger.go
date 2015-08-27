/*
Package llog is a logger for the main log.

Messages are written depending on the debug level.
Several helpers for consistent formating are provided.
Only 1 main log at one time is possible.
*/
package llog

import (
	"bytes"
	"log"
	"os"
)

// log levels
const (
	_       = iota
	MAIN    // Only main info
	ERROR   // display also Error messages
	WARNING // display also Warnings
	INFO    // display all Informations
	DEBUG   // display additional Debug messages
)

var ll lvlLogger

type lvlLogger struct {
	lvl int         // debug level
	out *os.File    // log file
	log *log.Logger // logger
}

// New creates a new main logger which stores
// the messages to filename. Messages are only written
// if the message log level is >= level.
func New(filename string, level int) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	l := log.New(f, "", log.LstdFlags|log.LUTC)
	ll = lvlLogger{lvl: level, out: f, log: l}
	return nil
}

// Close removes the main logger.
func Close() {
	ll.out.Close()
}

// SetLevel sets a new log level
func SetLevel(l int) {
	ll.lvl = l
}

// PrintMain writes a Maininfo msg to the log
func PrintMain(id, msg string) {
	ll.writeLog(MAIN, "%-15s :INFO : %s\n", id, msg)
}

// PrintInfo writes an Info msg to the log
func PrintInfo(id, msg string) {
	ll.writeLog(INFO, "%-15s :INFO : %s\n", id, msg)
}

// PrintWarning writes a Warning msg to the log
func PrintWarning(id, msg string) {
	ll.writeLog(WARNING, "%-15s :WARN : %s\n", id, msg)
}

// PrintError writes an Error msg to the log
func PrintError(id, msg string) {
	ll.writeLog(ERROR, "%-15s :ERROR: %s\n", id, msg)
}

// PrintDebug writes a Debug msg to the log
func PrintDebug(id, msg string) {
	ll.writeLog(DEBUG, "%-15s :DEBUG: %s\n", id, msg)
}

func (l lvlLogger) writeLog(level int, format string, msg ...interface{}) {
	if level <= l.lvl {
		l.log.Printf(format, msg...)
		l.out.Sync()
	}
}

func (l *lvlLogger) rotate(postfix string) error {
	var buf bytes.Buffer
	l.log.SetOutput(&buf)
	n := l.out.Name()
	l.out.Close()
	e := os.Rename(n, n+postfix)
	l.out, e = os.OpenFile(n, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if e != nil {
		return e
	}
	l.log.SetOutput(l.out)
	if buf.Len() > 0 {
		l.out.Write(buf.Bytes())
	}
	return nil
}

// Rotate moves the actual log file to the same file
// name plus an prefix of '.1'. If there is an existing
// file with that prefix, it is move to '.2'. This is done
// up to num itterations. The num + 1 file will be deleted.
func Rotate(num int) error {
	return ll.rotate(".1")
}
