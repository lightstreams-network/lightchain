package tracer

import (
	tmtLog "github.com/tendermint/tendermint/libs/log"
	"fmt"
)

type Config struct {
	ShouldTrace bool
	LogFilePath string
}

func NewConfig(shouldTrace bool, logFilePath string) Config {
	return Config{shouldTrace, logFilePath}
}

func (c Config) PrintWarning(logger tmtLog.Logger) {
	logger.Info("|--------")
	logger.Info("| Danger: Tracing enabled is not recommended in production!")
	logger.Info(fmt.Sprintf("| Tracing output is configured to be persisted at %v", c.LogFilePath))
	logger.Info("|--------")
}