package cmd

import (
	"strconv"

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
