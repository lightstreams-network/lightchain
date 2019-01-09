package main

import (
	"fmt"
	"strings"
	"gopkg.in/urfave/cli.v1"
)



// Major version component of the current release
const Major = "0"

// Minor version component of the current release
const Minor = "9"

// Fix version component of the current release
const Fix = "0"

var (
	// Version is the full version string
	Version = strings.Join([]string{Major, Minor, Fix}, ".");
)

func VersionCmd(ctx *cli.Context) error {
	fmt.Println("Version: ", Version)
	return nil
}
