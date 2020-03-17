package common

import (
	"fmt"
	"strconv"
)

func BuildFilename(filebase string, i, count int) string {
	pad := len(strconv.Itoa(count))
	return fmt.Sprintf(
		"%s_%0"+strconv.Itoa(pad)+"d-%d.json",
		filebase, i, count)
}
