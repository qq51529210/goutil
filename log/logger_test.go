package log

import (
	"log"
	"os"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func Test_Logger(t *testing.T) {
	lg := NewLogger(os.Stdout, DefaultHeader, "default")
	testLoggerDebug(t, lg)
	testLoggerInfo(t, lg)
	testLoggerWarn(t, lg)
	testLoggerError(t, lg)
	//
	lg = NewLogger(os.Stdout, FileNameHeader, "file name")
	testLoggerDebug(t, lg)
	testLoggerInfo(t, lg)
	testLoggerWarn(t, lg)
	testLoggerError(t, lg)
	//
	lg = NewLogger(os.Stdout, FilePathHeader, "file path")
	testLoggerDebug(t, lg)
	testLoggerInfo(t, lg)
	testLoggerWarn(t, lg)
	testLoggerError(t, lg)
	//
	SetLogger(lg)
	f := func() {
		ErrorfDepthTrace(0, "gogo", "%d", 1)
		log.Output(2, "123")
	}
	go f()
	f()
	//
	time.Sleep(time.Second * 3)
}

func testLoggerDebug(t *testing.T, lg *Logger) {
	lg.Debug(1, 2, 3)
	lg.Debugf("%d", 1)
	lg.DebugDepth(0, 1, 2, 3)
	lg.DebugfDepth(0, "%d", 1)
	lg.DebugTrace("trace", 1, 2)
	lg.DebugfTrace("trace", "%d", 1)
	lg.DebugDepthTrace(0, "trace", 1, 2)
	lg.DebugfDepthTrace(0, "trace", "%d", 1)
}

func testLoggerInfo(t *testing.T, lg *Logger) {
	lg.Info(1, 2, 3)
	lg.Infof("%d", 1)
	lg.InfoDepth(0, 1, 2, 3)
	lg.InfofDepth(0, "%d", 1)
	lg.InfoTrace("trace", 1, 2)
	lg.InfofTrace("trace", "%d", 1)
	lg.InfoDepthTrace(0, "trace", 1, 2)
	lg.InfofDepthTrace(0, "trace", "%d", 1)
}

func testLoggerWarn(t *testing.T, lg *Logger) {
	lg.Warn(1, 2, 3)
	lg.Warnf("%d", 1)
	lg.WarnDepth(0, 1, 2, 3)
	lg.WarnfDepth(0, "%d", 1)
	lg.WarnTrace("trace", 1, 2)
	lg.WarnfTrace("trace", "%d", 1)
	lg.WarnDepthTrace(0, "trace", 1, 2)
	lg.WarnfDepthTrace(0, "trace", "%d", 1)
}

func testLoggerError(t *testing.T, lg *Logger) {
	lg.Error(1, 2, 3)
	lg.Errorf("%d", 1)
	lg.ErrorDepth(0, 1, 2, 3)
	lg.ErrorfDepth(1, "%d", 1)
	lg.ErrorTrace("trace", 1, 2)
	lg.ErrorfTrace("trace", "%d", 1)
	lg.ErrorDepthTrace(0, "trace", 1, 2)
	lg.ErrorfDepthTrace(1, "trace", "%d", 1)
}

func Test_Recover(t *testing.T) {
	lg := NewLogger(os.Stdout, DefaultHeader, "default")
	defer func() {
		lg.Recover(recover())
	}()

	testRecover()
}

func testRecover() {
	testRecover1()
}

func testRecover1() {
	testRecover2()
}

func testRecover2() {
	panic("test recover")
}
