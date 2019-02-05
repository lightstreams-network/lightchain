package api

import "fmt"

func featureNotSupportedErr() error {
	return fmt.Errorf("feature not supported by lightchain yet")
}
