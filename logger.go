package logger

import (
	"fmt"
	"os"
	"strings"
)

// DefaultLogger default logger
var DefaultLogger = Namespace("")

// DefaultEnvironmentVariablePrefix default environment variable prefix
const DefaultEnvironmentVariablePrefix = "SEVERINO_LOGGER"

const (
	LevelError Level = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

type (
	Level uint
	// LoggerInterface ...
	LoggerInterface interface {
	}
	// LoggerInitInterface ...
	LoggerInitInterface interface {
		Init(namespace string, level Level)
	}
	// LoggerDebugInterface ...
	LoggerDebugInterface interface {
		Debug(msg string)
	}
	// LoggerInfoInterface ...
	LoggerInfoInterface interface {
		Info(msg string)
	}
	// LoggerWarnInterface ...
	LoggerWarnInterface interface {
		Warn(msg string)
	}
	// LoggerErrorInterface ...
	LoggerErrorInterface interface {
		Error(msg string)
	}
	// LoggerFatalInterface ...
	LoggerFatalInterface interface {
		Fatal(msg string)
	}

	Logger struct {
		Namespace string
		Level     Level
		Handlers  []LoggerInterface
	}
)

func getEnvVarLevel(namespace string) string {
	prefix := DefaultEnvironmentVariablePrefix
	if namespace != "" {
		prefix += "_"
		namespace = strings.ToUpper(namespace)
		namespace = strings.Replace(namespace, "-", "_", -1)
		namespace = strings.Replace(namespace, ".", "_", -1)
	}

	level := os.Getenv(prefix + namespace)
	if level == "" {
		level = os.Getenv(DefaultEnvironmentVariablePrefix)
	}

	return strings.ToLower(level)
}

func getLevelByString(level string) Level {
	if strings.EqualFold(level, "debug") {
		return LevelDebug
	} else if strings.EqualFold(level, "info") {
		return LevelInfo
	} else if strings.EqualFold(level, "warn") {
		return LevelWarn
	} else if strings.EqualFold(level, "error") {
		return LevelError
	} else {
		return LevelWarn
	}
}

// Namespace create a new logger namespace (new instance of logger)
func Namespace(namespace string) *Logger {
	logger := &Logger{
		Namespace: namespace,
	}

	logger.SetLevel(getLevelByString(getEnvVarLevel(namespace)))
	logger.AddHandler(&DefaultHandler{})

	return logger
}

func (logger *Logger) AddHandler(handler LoggerInterface) {
	logger.Handlers = append(logger.Handlers, handler)

	if initHandler, ok := handler.(LoggerInitInterface); ok {
		initHandler.Init(logger.Namespace, logger.Level)
	}
}

func (logger *Logger) SetLevel(level Level) {
	logger.Level = level

	for _, handler := range logger.Handlers {
		if initHandler, ok := handler.(LoggerInitInterface); ok {
			initHandler.Init(logger.Namespace, logger.Level)
		}
	}
}

func (logger *Logger) Debug(format string, v ...interface{}) {
	if logger.Level < LevelDebug {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if debugHandler, ok := handler.(LoggerDebugInterface); ok {
			debugHandler.Debug(msg)
		}
	}
}

func (logger *Logger) Info(format string, v ...interface{}) {
	if logger.Level < LevelInfo {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if infoHandler, ok := handler.(LoggerInfoInterface); ok {
			infoHandler.Info(msg)
		}
	}
}

func (logger *Logger) Warn(format string, v ...interface{}) {
	if logger.Level < LevelWarn {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if warnHandler, ok := handler.(LoggerWarnInterface); ok {
			warnHandler.Warn(msg)
		}
	}
}

func (logger *Logger) Error(format string, v ...interface{}) {
	if logger.Level < LevelError {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if errorHandler, ok := handler.(LoggerErrorInterface); ok {
			errorHandler.Error(msg)
		}
	}
}

func (logger *Logger) Fatal(format string, v ...interface{}) {
	if logger.Level < LevelError {
		return
	}

	msg := fmt.Sprintf(format, v...)
	for _, handler := range logger.Handlers {
		if fatalHandler, ok := handler.(LoggerFatalInterface); ok {
			fatalHandler.Fatal(msg)
		}
	}
	os.Exit(1)
}

func AddHandler(handler LoggerInterface) {
	DefaultLogger.AddHandler(handler)
}

func SetLevel(level Level) {
	DefaultLogger.SetLevel(level)
}

func Debug(format string, v ...interface{}) {
	DefaultLogger.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	DefaultLogger.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	DefaultLogger.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	DefaultLogger.Error(format, v...)
}

func Fatal(format string, v ...interface{}) {
	DefaultLogger.Fatal(format, v...)
}
