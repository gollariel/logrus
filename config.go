package logrus

import "github.com/kelseyhightower/envconfig"

// EnvPrefix environment prefix for log config
const EnvPrefix = "LOG"

// LogConfig contains name and level for log
type LogConfig struct {
	Name              string `default:"" split_words:"true"`
	Level             int8   `default:"0" split_words:"true"`
	DisableCaller     bool   `default:"false" split_words:"true"`
	DisableStacktrace bool   `default:"false" split_words:"true"`
}

// GetLogConfigFromEnv return log configs bases on environment variables
func GetLogConfigFromEnv() (*LogConfig, error) {
	c := new(LogConfig)
	err := envconfig.Process(EnvPrefix, c)
	return c, err
}
