package prometheus

import (
	"time"
)

type Config struct {
	enabled   bool
	http      HttpConfig
	namespace string
}

type HttpConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var (
	DefaultPrometheusAddr = ":26661"
	DefaultPrometheusNamespace = "lightchain"
)

func NewConfig(enabled bool, addr string, namespace string) Config {
	return Config{
		enabled: enabled,
		namespace: namespace,
		http: HttpConfig{
			Addr:         addr,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}
