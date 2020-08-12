package version

import (
	"fmt"
	"os"
	"runtime"
)

// set by build LD_FLAGS
var (
	version   string
	buildTime string
	gitCommit string
)

// Version info struct
type Version struct {
	Version   string `json:"version"`
	GitCommit string `json:"commit"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
	Os        string `json:"os"`
	Arch      string `json:"arch"`
}

// Get get Version instance
func Get() *Version {
	return &Version{
		Version:   version,
		GitCommit: gitCommit,
		BuildTime: buildTime,
		GoVersion: runtime.Version(),
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// ShowVersion show version info
func ShowVersion() {
	fmt.Println(Get())
	os.Exit(0)
}

// String return Version format string
func (v *Version) String() string {
	versionInfo := `Version:      ` + Get().Version + `
Git commit:   ` + Get().GitCommit + `
Go version:   ` + Get().GoVersion + `
Built time:   ` + Get().BuildTime + `
OS/Arch:      ` + Get().Os + "/" + Get().Arch

	return versionInfo
}
