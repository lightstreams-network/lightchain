package api

import "fmt"

const version = 1

func featureNotSupportedErr() error {
	return fmt.Errorf("feature not supported by lightchain yet")
}