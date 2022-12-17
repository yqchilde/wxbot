package log

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

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

type Formatter struct{}

func (s *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	msg := fmt.Sprintf("[%s] [%-5s] [%s:%d] %s\n", timestamp, level, log.callerFile, log.callerLine, entry.Message)
	return []byte(msg), nil
}

func init() {
	log.l.SetLevel(logrus.TraceLevel)
	log.l.SetOutput(os.Stdout)
	log.l.SetFormatter(&Formatter{})
}

func getCaller() {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return
	}
	log.callerFile = path.Join(path.Base(path.Dir(file)), path.Base(file))
	log.callerLine = line
}

func (log *logger) log(level logrus.Level, msg ...interface{}) {
	getCaller()
	log.l.Log(level, msg...)
}

func (log *logger) logf(level logrus.Level, format string, msg ...interface{}) {
	getCaller()
	log.l.Logf(level, format, msg...)
}

func Println(args ...interface{}) {
	log.log(logrus.InfoLevel, args...)
}

func Printf(format string, args ...interface{}) {
	log.logf(logrus.InfoLevel, format, args...)
}

func Debug(args ...interface{}) {
	log.log(logrus.DebugLevel, args...)
}

func Debugf(format string, args ...interface{}) {
	log.logf(logrus.DebugLevel, format, args...)
}

func Warn(args ...interface{}) {
	log.log(logrus.WarnLevel, args...)
}

func Warnf(format string, args ...interface{}) {
	log.logf(logrus.WarnLevel, format, args...)
}

func Error(args ...interface{}) {
	log.log(logrus.ErrorLevel, args...)
}

func Errorf(format string, args ...interface{}) {
	log.logf(logrus.ErrorLevel, format, args...)
}

func Fatal(args ...interface{}) {
	log.log(logrus.FatalLevel, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.logf(logrus.FatalLevel, format, args...)
}

func Panic(args ...interface{}) {
	log.log(logrus.PanicLevel, args...)
}

func Panicf(format string, args ...interface{}) {
	log.logf(logrus.PanicLevel, format, args...)
}

func Trace(args ...interface{}) {
	log.log(logrus.TraceLevel, args...)
}

func Tracef(format string, args ...interface{}) {
	log.logf(logrus.TraceLevel, format, args...)
}
