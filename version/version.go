package version

import (
	"fmt"
)

var (
	// Version semver for syncrepos (set through go build -ldflags)
	Version string
	// BuildUser is the user login who initiated the build (set through go build -ldflags)
	BuildUser string
	// GitTag is the result of git describe (set through go build -ldflags)
	GitTag string
	// BuildDate is the date of build (set through go build -ldflags)
	BuildDate string
)

// String displays all the version values
func String() string {
	res := ""
	res = res + fmt.Sprintf("Git Tag   : %s\n", GitTag)
	res = res + fmt.Sprintf("Build User: %s\n", BuildUser)
	res = res + fmt.Sprintf("Version   : %s\n", Version)
	res = res + fmt.Sprintf("BuildDate : %s\n", BuildDate)
	return res
}
