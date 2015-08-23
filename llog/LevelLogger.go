/*Package llog is a logger for the main log file.
Only 1 main log is possible.
Messages are written depending on the debug level.
Several helpers for consistent formating are provided.
*/
package llog

import (
	"log"
	"os"
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

// PrintInfo writes an Info msg to the log
func PrintInfo(level int, id, msg string) {
	ll.printf(level, "%-15s :INFO : %s\n", id, msg)
}

// PrintWarning writes a Warning msg to the log
func PrintWarning(level int, id, msg string) {
	ll.printf(level, "%-15s :WARN : %s\n", id, msg)
}

// PrintError writes an Error msg to the log
func PrintError(level int, id, msg string) {
	ll.printf(level, "%-15s :ERROR: %s\n", id, msg)
}

// PrintDebug writes an Debug msg to the log
func PrintDebug(level int, id, msg string) {
	ll.printf(level, "%-15s :DEBUG: %s\n", id, msg)
}

func (l lvlLogger) printf(level int, format string, msg ...interface{}) {
	if level <= l.lvl {
		l.log.Printf(format, msg...)
		l.out.Sync()
	}
}
