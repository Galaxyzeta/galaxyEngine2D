package lock

import (
	"sync/atomic"
	"time"
)

type SpinLock struct {
	elem int32
}

const locked = 1
const unlocked = 0

func (lock *SpinLock) Lock() {
	for !atomic.CompareAndSwapInt32(&lock.elem, unlocked, locked) {
		time.Sleep(time.Microsecond * 100)
	}
}

func (lock *SpinLock) Unlock() {
	for !atomic.CompareAndSwapInt32(&lock.elem, locked, unlocked) {
		time.Sleep(time.Microsecond * 100)
	}
}
