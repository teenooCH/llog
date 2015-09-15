package llog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"

	"github.com/xlab/handysort"
)

// log levels
const (
	MAIN    = iota // Only main info
	ERROR          // display also Error messages
	WARNING        // display also Warnings
	INFO           // display all Information
	DEBUG          // display additional Debug messages
)

// default format for the log strings
const deFormat string = "%-15s:%-5s: %s\n"

var ll lvlLogger

type lvlLogger struct {
	lvl int         // debug level
	out *os.File    // log file
	log *log.Logger // logger
	fmt string      // print format
}

// New creates a new logger which stores
// the messages to filename. Messages are only written
// if the message log level is >= level.
func New(filename string, level int) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	l := log.New(f, "", log.LstdFlags|log.LUTC)
	ll = lvlLogger{lvl: level, out: f, log: l, fmt: deFormat}
	return nil
}

// Close will clos the logger opened by New().
func Close() {
	ll.out.Close()
}

// SetLevel sets a new log level. A value < 0 will
// turn off the logging.
func SetLevel(l int) {
	ll.lvl = l
}

// SetFormat set a new formating string for the Print* functions.
// It has to include at least 3 %s expressions.
func SetFormat(f string) {
	ll.fmt = f
}

// PrintMain writes a Maininfo msg to the log
func PrintMain(id, msg string) {
	ll.writeLog(MAIN, ll.fmt, id, "INFO", msg)
}

// PrintInfo writes an Info msg to the log
func PrintInfo(id, msg string) {
	ll.writeLog(INFO, ll.fmt, id, "INFO", msg)
}

// PrintWarning writes a Warning msg to the log
func PrintWarning(id, msg string) {
	ll.writeLog(WARNING, ll.fmt, id, "WARN", msg)
}

// PrintError writes an Error msg to the log
func PrintError(id, msg string) {
	ll.writeLog(ERROR, ll.fmt, id, "ERROR", msg)
}

// PrintDebug writes a Debug msg to the log
func PrintDebug(id, msg string) {
	ll.writeLog(DEBUG, ll.fmt, id, "DEBUG", msg)
}

func (l lvlLogger) writeLog(level int, format string, msg ...interface{}) {
	if level <= l.lvl {
		l.log.Printf(format, msg...)
		l.out.Sync()
	}
}

// rotate moves the current file to filename+postfix
// and opens a new one with the old name.
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
// name plus a prefix of '.1'. If there is an existing
// file with that prefix, it is moved to '.2'. This is done
// up to num itterations. All files > num will be deleted.
func Rotate(num int) error {
	files, err := getFiles(ll.out.Name())
	if err != nil {
		return err
	}
	if files != nil {
		if len(files) >= num {
			for _, f := range files[num-1:] {
				os.Remove(f)
			}
			files = files[:num-1]
		}
		for i := len(files) - 1; i >= 0; i-- {
			os.Rename(files[i], fmt.Sprintf("%s.%d", ll.out.Name(), i+2))
		}
	}
	return ll.rotate(".1")
}

// getFiles returns a sorted list of files in the directory
// in which fname exists starting with fname (i.e. fname*).
func getFiles(fname string) ([]string, error) {
	dn := path.Dir(fname)
	bn := path.Base(fname)
	files, err := ioutil.ReadDir(dn)
	if err != nil {
		return nil, err
	}
	match := make([]os.FileInfo, len(files))
	count := 0
	for _, f := range files {
		if m, _ := path.Match(bn+".*", f.Name()); m {
			match[count] = f
			count++
		}
	}
	if count > 0 {
		res := make([]string, count)
		for i, f := range match[:count] {
			res[i] = fmt.Sprintf("%s/%s", dn, f.Name())
		}
		sort.Sort(handysort.Strings(res))
		return res, nil
	}
	return nil, nil
}
