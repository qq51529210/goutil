package log

import (
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	print(NewLogger(os.Stderr, DefaultHeader, "default", "a"))
	print(NewLogger(os.Stderr, FileNameHeader, "filename", "b"))
	print(NewLogger(os.Stderr, FilePathHeader, "filepath", "c"))
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
}
