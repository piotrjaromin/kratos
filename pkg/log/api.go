package log

// import (
// 	log "github.com/sirupsen/logrus"
// )

// AppInfo is used to provide basic configuration of logger
type AppInfo interface {
	Version() string
	LogLevel() string
}

type MiddlewareOpts struct {
	SkipEndpoints []string
}

// GetLogger created new instance of logger
func GetLogger(appInfo AppInfo) *Logger {

	logger := NewLogger()
	logger.SetLevelFromString(appInfo.LogLevel())
	// logger.SetFormatter(&log.JSONFormatter{})
	return logger
}
