package graphics

import "galaxyzeta.io/engine/infra/concurrency/lock"

type vboPool struct {
	allocatePtr int
	recyclePtr  int
	maxsize     int
	vboLeft     int
	pool        []uint32
	mu          lock.SpinLock
}

func InitVboPool(size int) {
	vboManager = &vboPool{
		allocatePtr: 0,
		recyclePtr:  0,
		maxsize:     size,
		vboLeft:     size,
		pool:        make([]uint32, 0, size),
		mu:          lock.SpinLock{},
	}
	for i := 0; i < vboManager.maxsize; i++ {
		vboManager.pool = append(vboManager.pool, GLNewVBO(1))
	}
}

func (vp *vboPool) Borrow() (ret uint32) {
	vp.mu.Lock()
	defer vp.mu.Unlock()
	ret = vp.pool[vp.allocatePtr]
	if vp.vboLeft <= 0 {
		panic("invalid operation")
	}
	vp.vboLeft--
	vp.allocatePtr++
	if vp.allocatePtr >= vp.maxsize {
		vp.allocatePtr = 0
	}
	return
}

func (vp *vboPool) Release(idx uint32) {
	vp.mu.Lock()
	defer vp.mu.Unlock()
	if vp.vboLeft >= vp.maxsize {
		panic("invalid operation")
	}
	vp.vboLeft++
	vp.recyclePtr++
	if vp.recyclePtr >= vp.maxsize {
		vp.recyclePtr = 0
	}
}
