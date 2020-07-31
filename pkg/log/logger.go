package log

import (
	"encoding/json"
	"io"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

// Logger extend logrus.Logger
type Logger struct {
	*logrus.Logger
}

// NewLogger return new logger instance
func NewLogger() *Logger {
	return &Logger{
		Logger: logrus.New(),
	}
}

func (l *Logger) Child() *Logger {
	child := NewLogger()
	child.SetFormatter(l.Formatter())
	return child
}

// Output return logger io.Writer
func (l *Logger) Output() io.Writer {
	return l.Out
}

// Level return logger level
func (l *Logger) Level() log.Lvl {
	return toEchoLevel(l.Logger.Level)
}

// SetLevel logger level
func (l *Logger) SetLevel(v log.Lvl) {
	l.Logger.Level = toLogrusLevel(v)
}

// SetLevel logger level
func (l *Logger) SetLevelFromString(level string) {
	l.Logger.Level = toLogrusLevelFromString(level)
}

// SetHeader logger header
// Managed by Logrus itself
// This function do nothing
func (l *Logger) SetHeader(h string) {
	// do nothing
}

// Formatter return logger formatter
func (l *Logger) Formatter() logrus.Formatter {
	return l.Logger.Formatter
}

// SetFormatter logger formatter
// Only support logrus formatter
func (l *Logger) SetFormatter(formatter logrus.Formatter) {
	l.Logger.Formatter = formatter
}

// Prefix return logger prefix
// This function do nothing
func (l *Logger) Prefix() string {
	return ""
}

// SetPrefix logger prefix
// This function do nothing
func (l *Logger) SetPrefix(p string) {
	// do nothing
}

// Printj output json of print level
func (l *Logger) Printj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Println(string(b))
}

// Debugj output message of debug level
func (l *Logger) Debugj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Debugln(string(b))
}

// Infoj output json of info level
func (l *Logger) Infoj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Infoln(string(b))
}

// Warnj output json of warn level
func (l *Logger) Warnj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Warnln(string(b))
}

// Errorj output json of error level
func (l *Logger) Errorj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Errorln(string(b))
}

// Fatalj output json of fatal level
func (l *Logger) Fatalj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Fatalln(string(b))
}

// Panicj output json of panic level
func (l *Logger) Panicj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Panicln(string(b))
}
