package unittest_test

import (
	"fmt"
	"testing"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/core"
	objs "galaxyzeta.io/engine/examples/testproj/userspace"
	"galaxyzeta.io/engine/parser"
)

type Wrapper struct {
	*Inner
}

type Inner struct {
	a int
}

func TestIO(t *testing.T) {
	i := &Inner{}
	v := Wrapper{
		i,
	}
	i.a = 3
	fmt.Println(v.a)
	fmt.Println(v.Inner.a)
	fmt.Println(i.a)
}

func TestFileParser(t *testing.T) {
	filePath := "../static/level/level.xml"
	t.Log(filePath)
	cfg := parser.ParseGameLevelFile(filePath)
	t.Log(cfg)
}

func TestAutowire(t *testing.T) {
	p := objs.TestPlayer{
		GameObject2D: base.NewGameObject2D("player"),
	}
	core.Inject(p)
	t.Log(p)
}
