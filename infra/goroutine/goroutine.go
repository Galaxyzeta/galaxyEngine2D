package goroutine

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// GoID gets current goroutine's ID.
// Don't use this unless you're in Debug mode.
func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
