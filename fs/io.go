package fs

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"golang.org/x/crypto/ssh/terminal"
)

func AskForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y(yes)/n(no)]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func PromptPassword(msg string) (string, error) {
	fmt.Println(msg)
	pwd, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(pwd)), nil
}