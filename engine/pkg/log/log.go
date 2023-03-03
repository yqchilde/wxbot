package log

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

type logger struct {
	l          *logrus.Logger
	callerFile string
	callerLine int
}

var log = &logger{
	l: logrus.New(),
}

func init() {
	log.l.SetLevel(logrus.TraceLevel)
	log.l.SetOutput(os.Stdout)
	log.l.SetReportCaller(true)
	log.l.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			if os.Getenv("DEBUG") == "true" || os.Getenv("DEBUG_LOG") == "true" {
				return "", fmt.Sprintf("[%s:%d] [GOID:%d]", log.callerFile, log.callerLine, getGoId())
			}
			return "", ""
		},
	})
}

func GetLogger() *logrus.Logger {
	return log.l
}

func getCaller() {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return
	}
	log.callerFile = path.Join(path.Base(path.Dir(file)), path.Base(file))
	log.callerLine = line
}

func Println(args ...interface{}) {
	getCaller()
	log.l.Println(args...)
}

func Printf(format string, args ...interface{}) {
	getCaller()
	log.l.Printf(format, args...)
}

func Debug(args ...interface{}) {
	getCaller()
	log.l.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	getCaller()
	log.l.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	getCaller()
	log.l.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	getCaller()
	log.l.Warnf(format, args...)
}

func Error(args ...interface{}) {
	getCaller()
	log.l.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	getCaller()
	log.l.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	getCaller()
	log.l.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	getCaller()
	log.l.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	getCaller()
	log.l.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	getCaller()
	log.l.Panicf(format, args...)
}

func Trace(args ...interface{}) {
	getCaller()
	log.l.Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	getCaller()
	log.l.Tracef(format, args...)
}

func getGoId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
