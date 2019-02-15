package log

import (
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"fmt"
)

func New(logFilePath string) (*zap.SugaredLogger, error) {
	if err := ioutil.WriteFile(logFilePath, []byte(""), os.ModePerm); err != nil {
		return nil, fmt.Errorf("error creating log file. %v", err)
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