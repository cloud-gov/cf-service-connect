package utils

import (
	"fmt"
	"os"
)

func Debugf(format string, args ...interface{}) {
	if os.Getenv("DEBUG") {
		fmt.Printf(format, args...)
	}
}
