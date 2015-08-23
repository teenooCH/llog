/*
Output expected similar to :

2015/08/21 21:34:12                 :INFO : Start of test
2015/08/21 21:34:12 Info(2)         :INFO : Line 1
2015/08/21 21:34:12 Info(1)         :INFO : Line 2
2015/08/21 21:34:12 Info(4)         :INFO : Line 3
2015/08/21 21:34:12 Warning(3)      :WARN : Line 1
2015/08/21 21:34:12 Warning(3)      :WARN : Line 2
2015/08/21 21:34:12 Error(2)        :ERROR: Line 1
2015/08/21 21:34:12 Error(4)        :ERROR: Line 2
2015/08/21 21:34:12 Debug(4)        :DEBUG: Line 1
2015/08/21 21:34:12                 :INFO : End of test

*/
package llog_test

import (
	"llog"
	"os"
	"testing"
)

var f = os.TempDir() + "/test.log"

// some short cuts for the print functions
var pi = llog.PrintInfo
var pw = llog.PrintWarning
var pe = llog.PrintError
var pd = llog.PrintDebug

func init() {
	os.Remove(f)
}
func Test_CreateLog(t *testing.T) {
	e := llog.New(f, 4)
	if e != nil {
		t.Error("failed to create log : " + e.Error())
		return
	}
	pi(1, "", "Start of test")
	llog.Close()
}
func Test_OpenLog(t *testing.T) {
	pe(1, "Error(1)", "Cannot be seen")
	e := llog.New(f, 4)
	if e != nil {
		t.Error("failed to open log : " + e.Error())
		return
	}
	pi(2, "Info(2)", "Line 1")
}
func Test_Print(t *testing.T) {
	pi(1, "Info(1)", "Line 2")
	pi(4, "Info(4)", "Line 3")
	pi(5, "Info(5)", "Should not be seen")
	pw(3, "Warning(3)", "Line 1")
	pw(3, "Warning(3)", "Line 2")
	pe(2, "Error(2)", "Line 1")
	pe(4, "Error(4)", "Line 2")
	pd(4, "Debug(4)", "Line 1")
	pi(1, "", "End of test")
}

func Test_CleanUp(t *testing.T) {
	llog.Close()
}
