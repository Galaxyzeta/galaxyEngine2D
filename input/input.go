package input

import (
	"galaxyzeta.io/engine/core"
)
import "galaxyzeta.io/engine/input/keys"

func IsKeyPressed(k keys.Key) (b bool) {
	return core.IsSetInputBuffer(keys.Action_KeyPress, k)
}

func IsKeyHeld(k keys.Key) (b bool) {
	return core.IsSetInputBuffer(keys.Action_KeyHold, k)
}

func IsKeyReleased(k keys.Key) (b bool) {
	return core.IsSetInputBuffer(keys.Action_KeyRelease, k)
}
