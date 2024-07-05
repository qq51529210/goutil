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
	fmt.Println(testError1())
	fmt.Println(testError2())

	fmt.Println(testError3().Log())
	fmt.Println(testError4().Log())
}

func testError1() error {
	return NewFileNameError(0, "3", 1, io.EOF)
}

func testError2() error {
	return NewFilePathError(0, "4", 2, io.EOF)
}

func testError3() *FileNameError[int] {
	return NewFileNameError(0, "3", 1, io.EOF)
}

func testError4() *FilePathError[int] {
	return NewFilePathError(0, "", 2, io.EOF)
}
