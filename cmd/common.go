package cmd

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/go-homedir"

	"github.com/pkg/errors"
)

func checkDataVersion(curVersion int32) int32 {
	if dataVersion != "" {
		dv, err := strconv.Atoi(dataVersion)
		checkError(errors.Wrap(err, "version invalid"))

		curVersion = int32(dv)
	}

	return curVersion
}

func saveConfigFile(s string) {
	home, err := homedir.Dir()
	checkError(errors.Wrap(err, "fail to get homedir"))

	filePath := filepath.Join(home, ".zkcmd.yaml")

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	checkError(err)
	defer f.Close()

	_, err = f.WriteString(s)
	checkError(err)
}
