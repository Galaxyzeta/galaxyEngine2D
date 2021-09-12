package unittest_test

import (
	"fmt"
	"testing"

	"galaxyzeta.io/engine/core"
	_ "galaxyzeta.io/engine/examples/testproj/userspace"
	"galaxyzeta.io/engine/sdk"
)

func init() {
	core.GlobalInitializer()
}

func TestExample(t *testing.T) {
	sdk.StartApplicationFromFile(fmt.Sprintf("%s/examples/testproj/static/level/level.xml", core.GetCwd()))
}
