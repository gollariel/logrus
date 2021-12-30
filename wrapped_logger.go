package logrus

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	skipCallers = 1
)

// WrappedLogger is a logger that is wrapped with a config file
// It takes a config to determine which level its running at
type WrappedLogger struct {
	*zap.SugaredLogger
}

// GetLoggerConfig return default zap config
func GetLoggerConfig() zap.Config {
	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapConfig
}

// NewWrappedLogger creates a new logger which wraps the given config
func NewWrappedLogger(name string, config zap.Config, args ...interface{}) (*WrappedLogger, error) {
	l, err := config.Build(zap.AddCallerSkip(skipCallers))
	if err != nil {
		return nil, err
	}

	if len(name) > 0 {
		l = l.Named(name)
	}
	sl := l.Sugar().With(args...)
	return &WrappedLogger{
		SugaredLogger: sl,
	}, nil
}

// WrapCore wrap zap core
func WrapCore(core zapcore.Core, options ...zap.Option) (*WrappedLogger, error) {
	options = append(options, zap.AddCallerSkip(skipCallers))
	return &WrappedLogger{
		SugaredLogger: zap.New(core, options...).Sugar(),
	}, nil
}

// Chain return new chain with local context
func (w *WrappedLogger) Chain() Chain {
	return &Entry{
		Logger: w.Desugar(),
	}
}

// With return new wrapped logger with default fields
func (w *WrappedLogger) With(v ...interface{}) *WrappedLogger {
	return &WrappedLogger{
		SugaredLogger: w.SugaredLogger.With(v...),
	}
}

// Print logs a message at Info level
func (w *WrappedLogger) Print(v ...interface{}) {
	w.SugaredLogger.Info(v...)
}

// Printf logs a message at Info level
func (w *WrappedLogger) Printf(format string, v ...interface{}) {
	w.SugaredLogger.Infof(format, v...)
}

// Println logs a message at Info level
func (w *WrappedLogger) Println(v ...interface{}) {
	w.SugaredLogger.Info(v...)
}

// Output add support nsq output
func (w *WrappedLogger) Output(calldepth int, s string) error {
	w.Desugar().Info(s, zap.Int("calldepth", calldepth))
	return nil
}
