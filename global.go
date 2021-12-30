package logrus

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	mu                 sync.RWMutex
	logConfig          *LogConfig
	globalLoggerStdout *WrappedLogger
	config             = GetLoggerConfig()
)

func init() {
	mu.Lock()
	defer mu.Unlock()

	cfg, err := GetLogConfigFromEnv()
	if err != nil {
		panic(fmt.Sprintf("Unable get log config, err: %s", err))
	}
	logConfig = cfg
	globalLoggerStdout = GetDefaultWrappedLogger(logConfig, config)
}

// GetDefaultWrappedLogger return default wrapper
func GetDefaultWrappedLogger(cfg *LogConfig, localCfg zap.Config) *WrappedLogger {
	localCfg.DisableCaller = cfg.DisableCaller
	localCfg.DisableStacktrace = cfg.DisableStacktrace
	localCfg.Level.SetLevel(zapcore.Level(cfg.Level))

	logger, err := NewWrappedLogger(cfg.Name, localCfg)
	if err != nil {
		panic(fmt.Sprintf("Unable get wrapped log, err: %s", err))
	}
	return logger
}

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// Logger return logger
func Logger() *WrappedLogger {
	mu.RLock()
	defer mu.RUnlock()
	return globalLoggerStdout
}

// GetChain return chain
func GetChain() Chain {
	mu.RLock()
	defer mu.RUnlock()
	return GetDefaultWrappedLogger(logConfig, config).Chain()
}

// SetLevel set login level
func SetLevel(level zapcore.Level) {
	config.Level.SetLevel(level)
}

// SetFormatter set formatter
func SetFormatter(encoding string, encoderConfig zapcore.EncoderConfig) {
	mu.Lock()
	defer mu.Unlock()
	config.Encoding = encoding
	config.EncoderConfig = encoderConfig
}

// SetOutput set output
func SetOutput(output string) {
	mu.Lock()
	defer mu.Unlock()
	config.OutputPaths = []string{output}
}

// SetErrorOutput set error output
func SetErrorOutput(errorOutput string) {
	mu.Lock()
	defer mu.Unlock()
	config.ErrorOutputPaths = []string{errorOutput}
}

// ApplyConfig apply new config to global logger
func ApplyConfig() (err error) {
	mu.Lock()
	defer mu.Unlock()
	globalLoggerStdout, err = NewWrappedLogger(logConfig.Name, config)
	return err
}

// With return logger with fields
func With(args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	globalLoggerStdout.SugaredLogger = Logger().SugaredLogger.With(args...)
}

// Debug prints a debug message
func Debug(message string, args ...interface{}) Chain {
	return GetChain().Debug(prepare(message, args...))
}

// Info logs a message with some additional context
func Info(message string, args ...interface{}) Chain {
	return GetChain().Info(prepare(message, args...))
}

// Warn logs a message with some additional context
func Warn(message string, args ...interface{}) Chain {
	return GetChain().Warn(prepare(message, args...))
}

// Warning logs a message with some additional context
func Warning(message string, args ...interface{}) Chain {
	return GetChain().Warn(prepare(message, args...))
}

// Error logs a message with some additional context
func Error(message string, args ...interface{}) Chain {
	return GetChain().Error(prepare(message, args...))
}

// Panic logs a message with some additional context, then panics
func Panic(message string, args ...interface{}) Chain {
	return GetChain().Panic(prepare(message, args...))
}

// Fatal logs a message to stderr and then calls an os.Exit
func Fatal(message string, args ...interface{}) Chain {
	return GetChain().Fatal(prepare(message, args...))
}

// Debugw prints a debug message
func Debugw(message string, keyValues ...interface{}) {
	Logger().Debugw(message, keyValues...)
}

// Infow logs a message with some additional context
func Infow(message string, keyValues ...interface{}) {
	Logger().Infow(message, keyValues...)
}

// Warnw logs a message with some additional context
func Warnw(message string, keyValues ...interface{}) {
	Logger().Warnw(message, keyValues...)
}

// Errorw logs a message with some additional context
func Errorw(message string, keyValues ...interface{}) {
	Logger().Errorw(message, keyValues...)
}

// Panicw logs a message with some additional context, then panics
func Panicw(message string, keyValues ...interface{}) {
	Logger().Panicw(message, keyValues...)
}

// Fatalw logs a message to stderr and then calls an os.Exit
func Fatalw(message string, keyValues ...interface{}) {
	Logger().Fatalw(message, keyValues...)
}

// Print logs a message at Info level
func Print(v ...interface{}) {
	Logger().Print(v...)
}

// Printf logs a message at Info level
func Printf(format string, v ...interface{}) {
	Logger().Printf(format, v...)
}

// Println logs a message at Info level
func Println(v ...interface{}) {
	Logger().Println(v...)
}

// WithField return logger with given field
func WithField(key string, value interface{}) Chain {
	return GetChain().WithAny(key, value)
}

// WithError return logger with given error
func WithError(err error) Chain {
	return GetChain().WithError(err)
}

// WithFields return logger with given fields
func WithFields(fields Fields) Chain {
	l := GetChain()
	for k, v := range fields {
		l = l.WithAny(k, v)
	}
	return l
}

func ChainWith(c Chain, keyvals ...interface{}) Chain {
	newChain := c.Copy()
	ln := len(keyvals)
	if ln%2 != 0 {
		return c
	}

	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}
		newChain.WithAny(key, keyvals[i+1])
	}

	return newChain
}
