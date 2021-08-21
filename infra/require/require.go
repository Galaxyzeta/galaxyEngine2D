package require

import (
	"fmt"
)

func commonFatal(expected interface{}, actual interface{}) {
	panic(fmt.Sprintf("FAIL: expected: %v, actual: %v", expected, actual))
}

func EqBool(expected bool, actual bool) {
	if expected != actual {
		commonFatal(expected, actual)
	}
}

func EqInt(expected int, actual int) {
	if expected != actual {
		commonFatal(expected, actual)
	}
}
