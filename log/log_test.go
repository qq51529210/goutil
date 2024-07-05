package log

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	print(NewLogger(os.Stderr, DefaultHeader, "default", "a", nil))
	print(NewLogger(os.Stderr, FileNameHeader, "filename", "b", nil))
	print(NewLogger(os.Stderr, FilePathHeader, "filepath", "c", nil))
}

func print(lg *Logger) {
	lg.Debug("1")
	lg.DebugTrace("b", 2)
	lg.DebugDepth(1, "2")
	lg.DebugDepthTrace(1, "c", 3)
	//
	lg.Debugf("%d", 4)
	lg.DebugfTrace("d", "%d", 5)
	lg.DebugfDepth(1, "%d", 6)
	lg.DebugfDepthTrace(1, "e", "%d", 7)
	//
	printPanic(lg)
}

func printPanic(lg *Logger) {
	defer func() {
		lg.Recover(recover())
	}()
	panic("test panice")
}

func TestError(t *testing.T) {

	fmt.Println(StatckError1())
	fmt.Println(StatckError2())

	fmt.Println(StatckError3().String())
	fmt.Println(StatckError4().String())
	fmt.Println(StatckError5().String())

	DefaultLogger.name = "[app]"
	DefaultLogger.module = "[test]"

	DefaultLogger.Debug(StatckError3().String())
	DefaultLogger.Debug(StatckError4().String())
	DefaultLogger.Debug(StatckError5().String())
}

func StatckError1() error {
	return NewFileNameError(0, "3", 1, io.EOF)
}

func StatckError2() error {
	return NewFilePathError(0, "4", 2, io.EOF)
}

func StatckError3() *StatckError[int] {
	return NewFileNameError(0, "3", 1, io.EOF)
}

func StatckError4() *StatckError[int] {
	return NewFilePathError(0, "", 2, io.EOF)
}

func StatckError5() *StatckError[int] {
	return &StatckError[int]{Err: io.EOF.Error()}
}
