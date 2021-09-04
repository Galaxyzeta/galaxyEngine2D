package unittest_test

import (
	"fmt"
	"testing"

	"galaxyzeta.io/engine/level"
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
	cfg := level.ParseGameLevelFile(filePath)
	t.Log(cfg)
}
