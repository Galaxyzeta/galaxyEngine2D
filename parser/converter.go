package parser

import (
	"strconv"
	"strings"

	"galaxyzeta.io/engine/linalg"
)

func MustParseNumericStringTuple(tuple string) linalg.Vector2f64 {
	splited := strings.Split(tuple, ",")
	if len(splited) != 2 {
		panic("should only contain one common")
	}
	f1, err := strconv.ParseFloat(splited[0], 64)
	if err != nil {
		panic(err)
	}
	f2, err := strconv.ParseFloat(splited[1], 64)
	if err != nil {
		panic(err)
	}
	return linalg.NewVector2f64(f1, f2)
}
