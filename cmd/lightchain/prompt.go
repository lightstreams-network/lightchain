package main

import (
	"strings"
	"golang.org/x/crypto/ssh/terminal"
)

func promptPassword() (string, error) {
	pwd, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(pwd)), nil
}