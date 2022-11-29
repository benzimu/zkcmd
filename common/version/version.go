package version

import (
	"fmt"
	"os"
	"runtime"
)

// set by build LD_FLAGS
var (
	version   string
	buildDate string
	gitCommit string
)

// Version info struct
type Version struct {
	Version   string `json:"version"`
	GitCommit string `json:"commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	Compiler  string `json:"compiler"`
}

// Get get Version instance
func Get() *Version {
	return &Version{
		Version:   version,
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Compiler:  runtime.Compiler,
	}
}

// ShowVersion show version info
func ShowVersion() {
	fmt.Println(Get())
	os.Exit(0)
}

// String return Version format string
func (v *Version) String() string {
	return fmt.Sprintf(`Version:      %s
Git commit:   %s
Go version:   %s
Built date:   %s
Platform:     %s`,
		Get().Version,
		Get().GitCommit,
		Get().GoVersion,
		Get().BuildDate,
		fmt.Sprintf("%s/%s", Get().OS, Get().Arch))
}
