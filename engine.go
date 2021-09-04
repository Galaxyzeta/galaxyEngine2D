package main

import (
	"galaxyzeta.io/engine/core"
	_ "galaxyzeta.io/engine/examples/testproj/userspace"
	"galaxyzeta.io/engine/sdk"
)

func init() {
	core.GlobalInitializer()
}

func main() {
	sdk.StartApplicationFromFile("examples/testproj/static/level/level.xml")
}
