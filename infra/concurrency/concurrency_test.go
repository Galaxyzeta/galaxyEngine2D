package cc_test

import (
	"sync"
	"sync/atomic"
	"testing"

	cc "galaxyzeta.io/engine/infra/concurrency"
)

func TestSynergyGate(t *testing.T) {
	thresh := 8
	sg := cc.NewSynergyGate(int64(thresh))
	// sg := cc.New(thresh)
	wg := sync.WaitGroup{}
	var indicator int64 = 0
	var check int64 = 0
	var fail bool = false
	for i := 0; i < thresh; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 10; j++ {
				atomic.AddInt64(&indicator, 1)
				sg.Wait()
				if atomic.CompareAndSwapInt64(&check, 0, 1) {
					// check once, the first one got the chance to check.
					if !atomic.CompareAndSwapInt64(&indicator, int64(thresh), 0) {
						// must be [thresh], else the synergy gate is not working.
						fail = true
						check = 0
						wg.Done()
						return
					}
				}
			}
			wg.Done()
		}()
	}
	if fail {
		t.Fatal("not workin")
	}
	wg.Wait()
}
