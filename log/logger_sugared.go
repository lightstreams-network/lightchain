package log

import (
	"go.uber.org/zap"
	"os"
	"fmt"
)

func New(logFilePath string) (*zap.SugaredLogger, error) {
	f, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.OutputPaths = []string{logFilePath}
	cfg.DisableStacktrace = true
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()

	log, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("error building log cfg. %v", err)
	}
	defer log.Sync()

	return log.Sugar(), nil
}