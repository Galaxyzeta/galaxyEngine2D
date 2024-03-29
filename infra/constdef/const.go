package constdef

import (
	"fmt"

	"galaxyzeta.io/engine/base"
)

var DefaultGameFunction = func(igobj2d base.IGameObject2D) {}

var AlwaysTrueFunction = func() bool {
	return true
}

var AlwaysFalseFunction = func() bool {
	return false
}

var ErrNotFound = fmt.Errorf("not found")
