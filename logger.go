package logger

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

// DefaultLogger default logger
var DefaultLogger = Namespace("")
var loggers = map[string]*Logger{}
var loggersLock sync.Mutex

// defaultEnvironmentVariablePrefix default environment variable prefix
var defaultEnvironmentVariablePrefix = "SEVERINO_LOGGER"

const (
	// LevelNone ...
	LevelNone Level = iota
	// LevelError ...
	LevelError
	// LevelWarn ...
	LevelWarn
	// LevelInfo ...
	LevelInfo
	// LevelDebug ...
	LevelDebug
)

type (
	// Level ...
	Level uint
	// Interface ...
	Interface interface {
	}
	// InitInterface ...
	InitInterface interface {
		Init(namespace string, level Level)
	}
	// DebugInterface ...
	DebugInterface interface {
		Debug(msg string)
	}
	// InfoInterface ...
	InfoInterface interface {
		Info(msg string)
	}
	// WarnInterface ...
	WarnInterface interface {
		Warn(msg string)
	}
	// ErrorInterface ...
	ErrorInterface interface {
		Error(msg string)
	}
	// FatalInterface ...
	FatalInterface interface {
		Fatal(msg string)
	}

	// Logger ...
	Logger struct {
		Namespace string
		Level     Level
		Handlers  []Interface
	}
)

func getEnvVarLevel(namespace string) string {
	prefix := defaultEnvironmentVariablePrefix
	if namespace != "" {
		prefix += "_"
		namespace = strings.ToUpper(namespace)
		namespace = strings.Replace(namespace, "-", "_", -1)
		namespace = strings.Replace(namespace, ".", "_", -1)
	}

	level := os.Getenv(prefix + namespace)
	if level == "" {
		level = os.Getenv(defaultEnvironmentVariablePrefix)
	}

	return strings.ToLower(level)
}

func setEnvironmentVariablePrefix(prefix string) error {
	loggersLock.Lock()
	defer loggersLock.Unlock()

	for namespace := range loggers {
		if namespace != "" {
			return errors.New("Cannot change prefix because some logs have already been use")
		}
	}

	delete(loggers, "")
	defaultEnvironmentVariablePrefix = prefix

	return nil
}

func SetDefaultEnvironmentVariablePrefix(prefix string) error {
	if err := setEnvironmentVariablePrefix(prefix); err != nil {
		return err
	}
	DefaultLogger = Namespace("")

	return nil
}

func GetDefaultEnvironmentVariablePrefix() string {
	return defaultEnvironmentVariablePrefix
}

// GetLevelByString ...
func GetLevelByString(level string) Level {
	if strings.EqualFold(level, "debug") {
		return LevelDebug
	} else if strings.EqualFold(level, "info") {
		return LevelInfo
	} else if strings.EqualFold(level, "warn") {
		return LevelWarn
	} else if strings.EqualFold(level, "error") {
		return LevelError
	} else if strings.EqualFold(level, "none") {
		return LevelNone
	} else {
		return LevelInfo
	}
}

// Namespace create a new logger namespace (new instance of logger)
func Namespace(namespace string) *Logger {
	loggersLock.Lock()
	defer loggersLock.Unlock()
	namespaceLower := strings.ToLower(namespace)
	if logger, ok := loggers[namespaceLower]; ok {
		return logger
	}

	logger := &Logger{
		Namespace: namespace,
	}

	logger.SetLevel(GetLevelByString(getEnvVarLevel(namespace)))
	logger.AddHandler(&DefaultHandler{})

	loggers[namespaceLower] = logger

	return logger
}

// AddHandler ...
func (logger *Logger) AddHandler(handler Interface) {
	logger.Handlers = append(logger.Handlers, handler)

	if initHandler, ok := handler.(InitInterface); ok {
		initHandler.Init(logger.Namespace, logger.Level)
	}
}

// SetLevel ...
func (logger *Logger) SetLevel(level Level) {
	logger.Level = level

	for _, handler := range logger.Handlers {
		if initHandler, ok := handler.(InitInterface); ok {
			initHandler.Init(logger.Namespace, logger.Level)
		}
	}
}

// Debug ...
func (logger *Logger) Debug(format string, v ...interface{}) {
	if logger.Level < LevelDebug {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if debugHandler, ok := handler.(DebugInterface); ok {
			debugHandler.Debug(msg)
		}
	}
}

// Info ...
func (logger *Logger) Info(format string, v ...interface{}) {
	if logger.Level < LevelInfo {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if infoHandler, ok := handler.(InfoInterface); ok {
			infoHandler.Info(msg)
		}
	}
}

// Warn ...
func (logger *Logger) Warn(format string, v ...interface{}) {
	if logger.Level < LevelWarn {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if warnHandler, ok := handler.(WarnInterface); ok {
			warnHandler.Warn(msg)
		}
	}
}

// Error ...
func (logger *Logger) Error(format string, v ...interface{}) {
	if logger.Level < LevelError {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if errorHandler, ok := handler.(ErrorInterface); ok {
			errorHandler.Error(msg)
		}
	}
}

// Fatal ...
func (logger *Logger) Fatal(format string, v ...interface{}) {
	if logger.Level < LevelError {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if fatalHandler, ok := handler.(FatalInterface); ok {
			fatalHandler.Fatal(msg)
		}
	}
	os.Exit(1)
}

// Write ...
func (logger *Logger) Write(b []byte) (int, error) {
	logger.Info("%s", strings.TrimRight(string(b), "\n"))
	return len(b), nil
}

// AddHandler ...
func AddHandler(handler Interface) {
	DefaultLogger.AddHandler(handler)
}

// SetLevel ...
func SetLevel(level Level) {
	DefaultLogger.SetLevel(level)
}

// Debug ...
func Debug(format string, v ...interface{}) {
	DefaultLogger.Debug(format, v...)
}

// Info ...
func Info(format string, v ...interface{}) {
	DefaultLogger.Info(format, v...)
}

// Warn ...
func Warn(format string, v ...interface{}) {
	DefaultLogger.Warn(format, v...)
}

// Error ...
func Error(format string, v ...interface{}) {
	DefaultLogger.Error(format, v...)
}

// Fatal ...
func Fatal(format string, v ...interface{}) {
	DefaultLogger.Fatal(format, v...)
}
