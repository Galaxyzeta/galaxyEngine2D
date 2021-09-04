package unittest_test

import (
	"testing"

	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/sdk"
)

func init() {
	core.GlobalInitializer()
}

func TestExample(t *testing.T) {
	sdk.StartApplicationFromFile("../static/level/level.xml")
}
