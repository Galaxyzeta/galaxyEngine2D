package test_test

import "testing"

type root struct {
	val int
}

type child struct {
	root
	val2 int
}

func TestPolymorphism(t *testing.T) {

}
