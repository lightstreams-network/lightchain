package version

import "strings"

// Major version component of the current release
const Major = "0"

// Minor version component of the current release
const Minor = "0"

// Fix version component of the current release
const Fix = "0"

var (
	// Version is the full version string
	Version = strings.Join([]string{Major, Minor, Fix}, ".");
)