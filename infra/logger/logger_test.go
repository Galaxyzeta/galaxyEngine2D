package logger

import (
	"testing"
	"time"
)

func TestAsyncLogger(t *testing.T) {
	k := New("test")
	for i := 0; i < 1024; i++ {
		k.Debugf("test logging async...")
	}
	time.Sleep(time.Second)
}

func TestSizeOfByteAndString(t *testing.T) {
	str := "asd"
	t.Log(len(str))
	t.Log(len([]byte(str)))

	str = "哈哈"
	t.Log(len(str))
	t.Log(len([]byte(str)))
}
