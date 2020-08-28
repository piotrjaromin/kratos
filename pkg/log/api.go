package log

import (
	logrus "github.com/sirupsen/logrus"
)

// AppInfo is used to provide basic configuration of logger
type AppInfo struct {
	Version  string
	LogLevel string
}

type MiddlewareOpts struct {
	SkipEndpoints []string
}

// GetLogger created new instance of logger
func GetLogger(appInfo AppInfo) *Logger {
	logger := NewLogger()
	logger.SetLevelFromString(appInfo.LogLevel)
	// logger.SetFormatter(&log.JSONFormatter{})
	return logger
}

// GetCLILogger created new instance of logger for CLI purposes
func GetCLILogger() *Logger {
	logger := NewLogger()
	logger.SetLevelFromString("DEBUG")
	logger.SetFormatter(&logrus.TextFormatter{})
	return logger
}
