package tracer

type Config struct {
	ShouldTrace bool
	LogFilePath string
}

func NewConfig(shouldTrace bool, logFilePath string) Config {
	return Config{shouldTrace, logFilePath}
}