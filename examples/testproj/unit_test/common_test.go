package unittest_test

import (
	"fmt"
	"testing"
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
