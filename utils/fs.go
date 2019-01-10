package utils

import (
	"io/ioutil"
	"runtime"
	"path/filepath"
	"os"
	"os/user"
	"errors"
	"fmt"
)

func ReadFileContent(genesisPath string) ([]byte, error) {
	genesisBlob, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		return nil, err
	}

	return genesisBlob, nil
}

func CreatePathIfNotExists(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		errors.New(fmt.Sprintf("Data folder err: %v", err))
	}
	return nil
}

func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Lightchain")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "Lightchain")
		} else {
			return filepath.Join(home, ".lightchain")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}