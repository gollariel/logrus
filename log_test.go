package logrus

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
	"testing"
)

func TestWrapper(t *testing.T) {
	Errorw("test")
	Errorw("test2")
	GetChain().WithString("test", "string").Error("result").Error("result")
	GetChain().WithAny("test", struct {
		Name    string
		private string
	}{
		Name:    "name",
		private: "private",
	}).Info("result")
	WithError(errors.New("test")).Info("test")
	GetChain().Log("info", "test")  //nolint:errcheck
	GetChain().Log("info", "test")  //nolint:errcheck
	GetChain().Log("error", "test") //nolint:errcheck
	GetChain().Output(3, "test")    //nolint:errcheck

	logger := GetChain()
	ChainWith(logger, "test", "value")
	logger.Log("info", "test") //nolint:errcheck
	ChainWith(logger, "test", "value")
	logger.Log("info", "test") //nolint:errcheck
	l1 := logger.WithString("abc", "1")
	l2 := l1.WithString("abc", "2")
	l3 := logger.WithString("abc", "3")
	l1.Info("test")
	l2.Info("test")
	l3.Info("test")

	wrapper := GetDefaultWrappedLogger(&LogConfig{}, GetLoggerConfig())
	wrapper.Info("error")
}

func TestLog(t *testing.T) {
	SetOutput("./log")
	SetLevel(1)
	err := ApplyConfig()
	if err != nil {
		t.Fatal(err)
	}
	log.Print()
	Infow("Test")  // not logged
	Debugw("Test") // not logged
	Errorw("Test")
	Errorw("Test")

	file, err := os.Open("./log")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.Remove("./log")
		if err != nil {
			t.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	var i int
	for scanner.Scan() {
		var d map[string]string
		err = json.Unmarshal(scanner.Bytes(), &d)
		if err != nil {
			t.Fatal(err)
		}
		i++
	}

	if i != 2 {
		t.Fatalf("error unexpected value: %d != %d", i, 2)
	}

	err = scanner.Err()
	if err != nil {
		t.Fatalf("error unexpected error: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("error unexpected error: %v", err)
	}
}

func exampleError() error {
	return errors.New("test error")
}

/*
BenchmarkSugarZap-12                     5854836               202.6 ns/op             7 B/op          0 allocs/op
BenchmarkLocalZap-12                     2470467               480.0 ns/op           360 B/op          5 allocs/op
BenchmarkFastlogWithChain-12              854090              1279 ns/op            1194 B/op         11 allocs/op
BenchmarkChain-12                        1181947              1007 ns/op            1188 B/op         11 allocs/op
BenchmarkChainWith-12                    1237754               936.9 ns/op           755 B/op         12 allocs/op
BenchmarkGoKit-12                         457401              2611 ns/op            1670 B/op         18 allocs/op
*/

func BenchmarkSugarZap(b *testing.B) {
	b.StopTimer()
	cfg := GetLoggerConfig()
	cfg.OutputPaths = []string{"/tmp/test.log"}
	l, err := NewWrappedLogger("wrapped_logger", cfg, "globalKey", "globalValue")
	if err != nil {
		b.Fatal(err)
	}
	err = exampleError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l.Errorw("error getting example", "err", err)
	}
}

func BenchmarkLocalZap(b *testing.B) {
	b.StopTimer()
	cfg := GetLoggerConfig()
	cfg.OutputPaths = []string{"/tmp/test.log"}
	l, err := NewWrappedLogger("wrapped_logger", cfg, "globalKey", "globalValue")
	if err != nil {
		b.Fatal(err)
	}
	err = exampleError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l.Chain().WithError(err).Error("error getting example")
	}
}

func BenchmarkFastlogWithChain(b *testing.B) {
	b.StopTimer()
	cfg := GetLoggerConfig()
	cfg.OutputPaths = []string{"/tmp/test.log"}
	l, err := NewWrappedLogger("wrapped_logger", cfg, "globalKey", "globalValue")
	if err != nil {
		b.Fatal(err)
	}
	err = exampleError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l.Chain().WithFields(Fields{
			"userId":    "test",
			"companyId": "test",
			"err":       err,
		}).Error("video.IssueCloudFrontCookieAccess: Video streaming is forbidden for company")
	}
}

func BenchmarkChain(b *testing.B) {
	b.StopTimer()
	cfg := GetLoggerConfig()
	cfg.OutputPaths = []string{"/tmp/test.log"}
	l, err := NewWrappedLogger("wrapped_logger", cfg, "globalKey", "globalValue")
	if err != nil {
		b.Fatal(err)
	}
	err = exampleError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l.Chain().WithError(err).WithAny("userId", "test").WithAny("companyId", "test").Info("video.IssueCloudFrontCookieAccess: Video streaming is forbidden for company")
	}
}

func BenchmarkChainWith(b *testing.B) {
	b.StopTimer()
	cfg := GetLoggerConfig()
	cfg.OutputPaths = []string{"/tmp/test.log"}
	l, err := NewWrappedLogger("wrapped_logger", cfg, "globalKey", "globalValue")
	if err != nil {
		b.Fatal(err)
	}
	err = exampleError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		logger := ChainWith(l.Chain(), "userId", "test")
		logger = ChainWith(logger, "companyId", "test")
		logger = ChainWith(logger, "err", err)
		logger.Log("video.IssueCloudFrontCookieAccess: Video streaming is forbidden for company") //nolint:errcheck
	}
}

func BenchmarkGoKit(b *testing.B) {
	b.StopTimer()
	cfg := GetLoggerConfig()
	cfg.OutputPaths = []string{"/tmp/test.log"}
	l, err := NewWrappedLogger("wrapped_logger", cfg, "globalKey", "globalValue")
	if err != nil {
		b.Fatal(err)
	}
	err = exampleError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l.Chain().WithError(err).WithAny("userId", "test").WithAny("companyId", "test").Log("video.IssueCloudFrontCookieAccess: Video streaming is forbidden for company") //nolint:errcheck
	}
}

func BenchmarkWithError(b *testing.B) {
	b.StopTimer()
	cfg := GetLoggerConfig()
	cfg.OutputPaths = []string{"/tmp/test.log"}
	l, err := NewWrappedLogger("wrapped_logger", cfg, "globalKey", "globalValue")
	if err != nil {
		b.Fatal(err)
	}
	chain := l.Chain()
	err = exampleError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		chain.WithError(err).WithFields(Fields{
			"test":  "value",
			"test2": "value2",
		}).Error("Something happen") //nolint:errcheck
	}
}

func BenchmarkWithFields(b *testing.B) {
	b.StopTimer()
	cfg := GetLoggerConfig()
	cfg.OutputPaths = []string{"/tmp/test.log"}
	l, err := NewWrappedLogger("wrapped_logger", cfg, "globalKey", "globalValue")
	if err != nil {
		b.Fatal(err)
	}
	chain := l.Chain()
	err = exampleError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		chain.WithFields(Fields{
			"err":   err,
			"test":  "value",
			"test2": "value2",
		}).Error("Something happen") //nolint:errcheck
	}
}
