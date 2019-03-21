package main

import (
	"bufio"
	"os"
	"strings"
)

func promptPassword() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	pw, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(pw), nil
}