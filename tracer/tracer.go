package tracer

import (
	"go.uber.org/zap"
	"github.com/lightstreams-network/lightchain/log"
)

type Tracer struct {
	shouldTrace bool
	Logger      zap.SugaredLogger
}

func NewTracer(cfg Config) (Tracer, error) {
	logger, err := log.New(cfg.LogFilePath)
	if err != nil {
		return Tracer{}, err
	}

	logger = logger.With("engine", "tracer")
	return Tracer{
		shouldTrace: cfg.ShouldTrace,
		Logger:      *logger,
	}, nil
}

func (t Tracer) Assert(assertion func(tracer Tracer)) {
	if ! t.shouldTrace {
		return
	}

	assertion(t)
}
