package logrus

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Entry contains values of given chain
type Entry struct {
	Logger    *zap.Logger
	keyValues []zap.Field
	mu        sync.RWMutex
}

func (c *Entry) getFields() []zap.Field {
	c.mu.RLock()
	defer c.mu.RUnlock()
	fields := make([]zap.Field, len(c.keyValues))
	copy(fields, c.keyValues)
	return fields
}

func (c *Entry) copyWith(fields ...zap.Field) *Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()

	newEntry := &Entry{
		Logger: c.Logger,
	}
	newEntry.keyValues = make([]zap.Field, len(c.keyValues))
	copy(newEntry.keyValues, c.keyValues)
	newEntry.keyValues = append(newEntry.keyValues, fields...)
	return newEntry
}

// Debug logs a message at DebugLevel
func (c *Entry) Debug(msg string, args ...interface{}) Chain {
	c.Logger.Debug(prepare(msg, args...), c.getFields()...)
	return c
}

// Info logs a message at InfoLevel
func (c *Entry) Info(msg string, args ...interface{}) Chain {
	c.Logger.Info(prepare(msg, args...), c.getFields()...)
	return c
}

// Warn logs a message at WarnLevel
func (c *Entry) Warn(msg string, args ...interface{}) Chain {
	c.Logger.Warn(prepare(msg, args...), c.getFields()...)
	return c
}

// Error logs a message at ErrorLevel
func (c *Entry) Error(msg string, args ...interface{}) Chain {
	c.Logger.Error(prepare(msg, args...), c.getFields()...)
	return c
}

// Fatal logs a message at FatalLevel
func (c *Entry) Fatal(msg string, args ...interface{}) Chain {
	c.Logger.Fatal(prepare(msg, args...), c.getFields()...)
	return c
}

// Panic logs a message at PanicLevel
func (c *Entry) Panic(msg string, args ...interface{}) Chain {
	c.Logger.Panic(prepare(msg, args...), c.getFields()...)
	return c
}

// Warning logs a message at WarnLevel
func (c *Entry) Warning(msg string, args ...interface{}) Chain {
	c.Logger.Warn(prepare(msg, args...), c.getFields()...)
	return c
}

func prepare(msg string, args ...interface{}) string {
	if len(args) > 0 {
		if msg == "" {
			msg = fmt.Sprint(args...)
		} else {
			msg = fmt.Sprintf(msg, args...)
		}
	}
	return msg
}

// WithBinary add a field that carries an opaque binary blob.
func (c *Entry) WithBinary(key string, val []byte) Chain {
	return c.copyWith(zap.Binary(key, val))
}

// WithBool add a field that carries a bool.
func (c *Entry) WithBool(key string, val bool) Chain {
	return c.copyWith(zap.Bool(key, val))
}

// WithByteString add a field that carries UTF-8 encoded text as a []byte.
func (c *Entry) WithByteString(key string, val []byte) Chain {
	return c.copyWith(zap.ByteString(key, val))
}

// WithComplex128 add a field that carries a complex number.
func (c *Entry) WithComplex128(key string, val complex128) Chain {
	return c.copyWith(zap.Complex128(key, val))
}

// WithComplex64 add a field that carries a complex number
func (c *Entry) WithComplex64(key string, val complex64) Chain {
	return c.copyWith(zap.Complex64(key, val))
}

// WithFloat64 add a field that carries a float64
func (c *Entry) WithFloat64(key string, val float64) Chain {
	return c.copyWith(zap.Float64(key, val))
}

// WithFloat32 add a field that carries a float32
func (c *Entry) WithFloat32(key string, val float32) Chain {
	return c.copyWith(zap.Float32(key, val))
}

// WithInt add a field with the given key and value.
func (c *Entry) WithInt(key string, val int) Chain {
	return c.copyWith(zap.Int(key, val))
}

// WithInt64 add a field with the given key and value.
func (c *Entry) WithInt64(key string, val int64) Chain {
	return c.copyWith(zap.Int64(key, val))
}

// WithInt32 add a field with the given key and value.
func (c *Entry) WithInt32(key string, val int32) Chain {
	return c.copyWith(zap.Int32(key, val))
}

// WithInt16 add a field with the given key and value.
func (c *Entry) WithInt16(key string, val int16) Chain {
	return c.copyWith(zap.Int16(key, val))
}

// WithInt8 add a field with the given key and value.
func (c *Entry) WithInt8(key string, val int8) Chain {
	return c.copyWith(zap.Int8(key, val))
}

// WithString add a field with the given key and value.
func (c *Entry) WithString(key string, val string) Chain {
	return c.copyWith(zap.String(key, val))
}

// WithError add a field with the given error value.
func (c *Entry) WithError(val error) Chain {
	return c.copyWith(zap.Error(val))
}

// WithNamedError add a field with the given key and value.
func (c *Entry) WithNamedError(key string, err error) Chain {
	return c.copyWith(zap.NamedError(key, err))
}

// WithUint add a field with the given key and value.
func (c *Entry) WithUint(key string, val uint) Chain {
	return c.copyWith(zap.Uint(key, val))
}

// WithUint64 add a field with the given key and value.
func (c *Entry) WithUint64(key string, val uint64) Chain {
	return c.copyWith(zap.Uint64(key, val))
}

// WithUint32 add a field with the given key and value.
func (c *Entry) WithUint32(key string, val uint32) Chain {
	return c.copyWith(zap.Uint32(key, val))
}

// WithUint16 add a field with the given key and value.
func (c *Entry) WithUint16(key string, val uint16) Chain {
	return c.copyWith(zap.Uint16(key, val))
}

// WithUint8 add a field with the given key and value.
func (c *Entry) WithUint8(key string, val uint8) Chain {
	return c.copyWith(zap.Uint8(key, val))
}

// WithUintptr add a field with the given key and value.
func (c *Entry) WithUintptr(key string, val uintptr) Chain {
	return c.copyWith(zap.Uintptr(key, val))
}

// WithReflect constructs a field with the given key and value.
func (c *Entry) WithReflect(key string, val interface{}) Chain {
	return c.copyWith(zap.Reflect(key, val))
}

// WithNamespace creates a named, isolated scope within the logger's context
func (c *Entry) WithNamespace(key string) Chain {
	return c.copyWith(zap.Namespace(key))
}

// WithStringer add a field with the given key and value.
func (c *Entry) WithStringer(key string, val fmt.Stringer) Chain {
	return c.copyWith(zap.Stringer(key, val))
}

// WithTime add a field with the given key and value.
func (c *Entry) WithTime(key string, val time.Time) Chain {
	return c.copyWith(zap.Time(key, val))
}

// WithStack constructs a field that stores a stacktrace
func (c *Entry) WithStack(key string) Chain {
	return c.copyWith(zap.Stack(key))
}

// WithDuration add a field with the given key and value.
func (c *Entry) WithDuration(key string, val time.Duration) Chain {
	return c.copyWith(zap.Duration(key, val))
}

// WithAny takes a key and an arbitrary value and chooses the best way to represent  them as a field
func (c *Entry) WithAny(key string, val interface{}) Chain {
	return c.copyWith(zap.Any(key, val))
}

// Print logs a message at PanicLevel
func (c *Entry) Print(v ...interface{}) {
	c.Logger.Sugar().Info(v...)
}

// Printf logs a message at PanicLevel
func (c *Entry) Printf(format string, v ...interface{}) {
	c.Logger.Sugar().Infof(format, v...)
}

// Println logs a message at PanicLevel
func (c *Entry) Println(v ...interface{}) {
	c.Logger.Sugar().Info(v...)
}

// WithField return logger with given field
func (c *Entry) WithField(key string, value interface{}) Chain {
	return c.WithAny(key, value)
}

// WithFields return logger with given fields
func (c *Entry) WithFields(fields Fields) Chain {
	zapFields := make([]zap.Field, len(fields))
	var i int
	for k, v := range fields {
		zapFields[i] = zap.Any(k, v)
		i++
	}
	return c.copyWith(zapFields...)
}

// Log logs as Info
func (c *Entry) Log(keyvals ...interface{}) error {
	globalKeyVals := make([]interface{}, len(c.keyValues))
	for i, v := range c.getFields() {
		globalKeyVals[i] = v
	}

	ln := len(keyvals)
	if ln > 1 && ln%2 == 0 {
		msg, ok := keyvals[1].(string)
		if ok {
			switch keyvals[0] {
			case "warn":
				c.Logger.Sugar().Warnw(msg, append(keyvals[2:], globalKeyVals...)...)
				return nil
			case "info":
				c.Logger.Sugar().Infow(msg, append(keyvals[2:], globalKeyVals...)...)
				return nil
			case "error":
				c.Logger.Sugar().Errorw(msg, append(keyvals[2:], globalKeyVals...)...)
				return nil
			}
		}
		c.Logger.Sugar().Infow("", append(keyvals, globalKeyVals...)...)
		return nil
	}
	c.Logger.Sugar().Info(append(keyvals, globalKeyVals...)...)
	return nil
}

// Output Add support nsq output
func (c *Entry) Output(calldepth int, s string) error {
	keyvals := make([]zap.Field, len(c.keyValues))
	copy(keyvals, c.getFields())
	keyvals = append(keyvals, zap.Int("calldepth", calldepth))
	c.Logger.Info(s, keyvals...)
	return nil
}

// Copy copy Chain
func (c *Entry) Copy() Chain {
	c.mu.Lock()
	defer c.mu.Unlock()

	newEntry := &Entry{
		Logger: c.Logger,
	}
	newEntry.keyValues = make([]zap.Field, len(c.keyValues))
	copy(newEntry.keyValues, c.keyValues)
	return newEntry
}
