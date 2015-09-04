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
	_       = iota
	MAIN    // Only main info
	ERROR   // display also Error messages
	WARNING // display also Warnings
	INFO    // display all Information
	DEBUG   // display additional Debug messages
)

const leFormat string = "%-15s :%-5s: %s\n"

var ll lvlLogger

type lvlLogger struct {
	lvl int         // debug level
	out *os.File    // log file
	log *log.Logger // logger
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
	ll = lvlLogger{lvl: level, out: f, log: l}
	return nil
}

// Close closes the file opened by New().
func Close() {
	ll.out.Close()
}

// SetLevel sets a new log level
func SetLevel(l int) {
	ll.lvl = l
}

// PrintMain writes a Maininfo msg to the log
func PrintMain(id, msg string) {
	ll.writeLog(MAIN, leFormat, id, "INFO", msg)
}

// PrintInfo writes an Info msg to the log
func PrintInfo(id, msg string) {
	ll.writeLog(INFO, leFormat, id, "INFO", msg)
}

// PrintWarning writes a Warning msg to the log
func PrintWarning(id, msg string) {
	ll.writeLog(WARNING, leFormat, id, "WARN", msg)
}

// PrintError writes an Error msg to the log
func PrintError(id, msg string) {
	ll.writeLog(ERROR, leFormat, id, "ERROR", msg)
}

// PrintDebug writes a Debug msg to the log
func PrintDebug(id, msg string) {
	ll.writeLog(DEBUG, leFormat, id, "DEBUG", msg)
}

func (l lvlLogger) writeLog(level int, format string, msg ...interface{}) {
	if level <= l.lvl {
		l.log.Printf(format, msg...)
		l.out.Sync()
	}
}

// rotate moves the current file to filename.postfix
// and opens a new one.
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
// file with that prefix, it is moved to '.2'. This is done
// up to num itterations. The num + 1 file will be deleted.
func Rotate(num int) error {
	dn := path.Dir(ll.out.Name())
	bn := path.Base(ll.out.Name())
	files, err := getFiles(dn, bn)
	if err != nil {
		return err
	}
	if files != nil {
		if len(files) >= num {
			for _, f := range files[num-1 : len(files)] {
				os.Remove(f)
			}
			files = files[0 : num-1]
		}
		for i := len(files) - 1; i >= 0; i-- {
			os.Rename(files[i], fmt.Sprintf("%s.%d", ll.out.Name(), i+2))
		}
	}
	return ll.rotate(".1")
}

// getFiles returns a sorted list of files in directory fpath
// starting with fname (i.e. path/fname*).
func getFiles(fpath, fname string) ([]string, error) {
	files, err := ioutil.ReadDir(fpath)
	if err != nil {
		return nil, err
	}
	match := make([]os.FileInfo, len(files))
	count := 0
	for _, f := range files {
		if m, _ := path.Match(fname+".*", f.Name()); m {
			match[count] = f
			count++
		}
	}
	if count > 0 {
		res := make([]string, count)
		for i, f := range match[0:count] {
			res[i] = fmt.Sprintf("%s/%s", fpath, f.Name())
		}
		sort.Sort(handysort.Strings(res))
		return res, nil
	}
	return nil, nil
}
