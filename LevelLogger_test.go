/*
Output expected similar to :

2015/09/15 21:31:56                :INFO : Start of test
2015/09/15 21:31:56 Info           :INFO : Line 1
2015/09/15 21:31:56 Info           :INFO : Line 2
2015/09/15 21:31:56 Info           :INFO : Line 3
2015/09/15 21:31:56 Warning        :WARN : Line 1
2015/09/15 21:31:56 Warning        :WARN : Line 2
2015/09/15 21:31:56 Error          :ERROR: Line 1
2015/09/15 21:31:56 Error          :ERROR: Line 2
2015/09/15 21:31:56 Debug          :DEBUG: Line 1
2015/09/15 21:31:56 Debug     :DEBUG     : Line 2
2015/09/15 21:31:56           :INFO      : End of test

*/
package llog_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/teenooCH/llog"
)

var f = os.TempDir() + "/test.log"
var fr = os.TempDir() + "/testrotate.log"

// some short cuts for the print functions
var pm = llog.PrintMain
var pi = llog.PrintInfo
var pw = llog.PrintWarning
var pe = llog.PrintError
var pd = llog.PrintDebug

func init() {
	os.Remove(f)
}

func Test_CreateLogFailure(t *testing.T) {
	e := llog.New("/foo/bar", llog.INFO)
	if e == nil {
		t.Error("failed to provoke creation error!")
		return
	}
}

func Test_CreateLog(t *testing.T) {
	e := llog.New(f, llog.INFO)
	if e != nil {
		t.Error("failed to create log : " + e.Error())
		return
	}
	pm("", "Start of test")
	llog.Close()
}

func Test_OpenLog(t *testing.T) {
	pe("Error", "Cannot be seen")
	e := llog.New(f, llog.INFO)
	if e != nil {
		t.Error("failed to open log : " + e.Error())
		return
	}
	pi("Info", "Line 1")
}

func Test_Print(t *testing.T) {
	pi("Info", "Line 2")
	pi("Info", "Line 3")
	pw("Warning", "Line 1")
	pw("Warning", "Line 2")
	pe("Error", "Line 1")
	pe("Error", "Line 2")
	pd("Debug", "Line not shown") // not shown in this test
}

func Test_SetDebugLevel(t *testing.T) {
	llog.SetLevel(llog.DEBUG)
}

func Test_PrintDebug(t *testing.T) {
	pd("Debug", "Line 1")
}

func Test_SetFormat(t *testing.T) {
	llog.SetFormat("%-10s:%-10s: %s\n")
	pd("Debug", "Line 2")
}

func Test_EndLogging(t *testing.T) {
	pm("", "End of test")
	llog.Close()
}

func Test_Rotate(t *testing.T) {
	e := llog.New(fr, llog.INFO)
	if e != nil {
		t.Error("failed to open log : " + e.Error())
		return
	}
	pi("Info", "Line 1 pre rotate")
	pi("Info", "Line 2 pre rotate")
	e = llog.Rotate(3)
	if e != nil {
		t.Error("failed to rotate log : " + e.Error())
		return
	}
	pi("Info", "Line 1 after rotate")
	pi("Info", "Line 2 after rotate")
}

func Example() {
	e := llog.New("/foo/bar", llog.INFO)
	if e != nil {
		fmt.Println("failed to create log : " + e.Error())
		return
	}
	llog.PrintInfo("My ID", "Test of llog")
	llog.Close()
}
