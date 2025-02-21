package log

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestPrint(t *testing.T) {
	lg := NewLogger(os.Stderr, "AppName", "API")
	print(lg)
	lg.FormatStackHeader = FileNameHeader
	print(lg)
	lg.FormatStackHeader = FilePathHeader
	print(lg)
}

func print(lg *Logger) {
	lg.Debug(-1, "", 0, "1")
	lg.Debug(-1, "", 0, "%d", 2)
	lg.Debug(-1, "t1", 0, "3")
	lg.Debug(-1, "t2", 0, "%d", 4)
	lg.Debug(0, "t3", time.Microsecond, "4")
	lg.Debug(0, "t4", time.Microsecond, "%d", 5)
	//
	printPanic(lg)
	lg.Debug(-1, "", 0, "--------------------------------------------")
}

func printPanic(lg *Logger) {
	defer func() {
		lg.Recover(recover())
	}()
	panic("test panic")
}

func TestError(t *testing.T) {

	fmt.Println(StatckError1())
	fmt.Println(StatckError2())

	fmt.Println(StatckError3().String())
	fmt.Println(StatckError4().String())
	fmt.Println(StatckError5().String())

	DefaultLogger.name = "[app]"
	DefaultLogger.module = "[test]"

	DefaultLogger.Debug(0, "TestError", 0, StatckError3().String())
	DefaultLogger.Debug(0, "TestError", 0, StatckError4().String())
	DefaultLogger.Debug(0, "TestError", 0, StatckError5().String())
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

func Test_Log_IntRightAlign(t *testing.T) {
	l := new(Log)
	l.IntLeftAlign(123456, 9)
	fmt.Println(string(l.b))
	l.Reset()
	l.IntRightAlign(123456, 8)
	fmt.Println(string(l.b))
}
