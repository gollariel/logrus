package logrus

import (
	"fmt"
	"time"
)

// Log describe log wrapper interface
type Log interface {
	Debugw(msg string, keyValues ...interface{})
	Infow(msg string, keyValues ...interface{})
	Warnw(msg string, keyValues ...interface{})
	Errorw(msg string, keyValues ...interface{})
	Panicw(msg string, keyValues ...interface{})
	Fatalw(msg string, keyValues ...interface{})

	Chain() Chain
}

// Chain context provide functionality for log with predefined values
type Chain interface {
	Copy() Chain
	Debug(msg string, args ...interface{}) Chain
	Info(msg string, args ...interface{}) Chain
	Warn(msg string, args ...interface{}) Chain
	Warning(msg string, args ...interface{}) Chain
	Error(msg string, args ...interface{}) Chain
	Fatal(msg string, args ...interface{}) Chain
	Panic(msg string, args ...interface{}) Chain

	WithBinary(key string, val []byte) Chain
	WithBool(key string, val bool) Chain
	WithByteString(key string, val []byte) Chain
	WithComplex128(key string, val complex128) Chain
	WithComplex64(key string, val complex64) Chain
	WithFloat64(key string, val float64) Chain
	WithFloat32(key string, val float32) Chain
	WithInt(key string, val int) Chain
	WithInt64(key string, val int64) Chain
	WithInt32(key string, val int32) Chain
	WithInt16(key string, val int16) Chain
	WithInt8(key string, val int8) Chain
	WithString(key string, val string) Chain
	WithError(val error) Chain
	WithNamedError(key string, err error) Chain
	WithUint(key string, val uint) Chain
	WithUint64(key string, val uint64) Chain
	WithUint32(key string, val uint32) Chain
	WithUint16(key string, val uint16) Chain
	WithUint8(key string, val uint8) Chain
	WithUintptr(key string, val uintptr) Chain
	WithReflect(key string, val interface{}) Chain
	WithNamespace(key string) Chain
	WithStringer(key string, val fmt.Stringer) Chain
	WithTime(key string, val time.Time) Chain
	WithStack(key string) Chain
	WithDuration(key string, val time.Duration) Chain
	WithAny(key string, val interface{}) Chain
	WithField(key string, value interface{}) Chain
	WithFields(fields Fields) Chain

	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})

	Log(keyvals ...interface{}) error
	Output(calldepth int, s string) error
}
