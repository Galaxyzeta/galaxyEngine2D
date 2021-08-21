package graphics

import (
	"container/list"
	"time"

	"galaxyzeta.io/engine/infra/concurrency/lock"
)

type vboPool struct {
	maxsize    int
	vboLeft    int
	vacantList *list.List
	inuseList  *list.List
	inuseIndex map[uint32]*list.Element //this can be used to look up for a specific node in list quickly.
	mu         lock.SpinLock
}

func InitVboPool(size int) {
	vboManager = &vboPool{
		maxsize:    size,
		vboLeft:    size,
		vacantList: list.New(),
		inuseList:  list.New(),
		inuseIndex: make(map[uint32]*list.Element),
		mu:         lock.SpinLock{},
	}

	for i := 0; i < vboManager.maxsize; i++ {
		newBuffer := GLNewVBO(1)
		vboManager.vacantList.PushBack(newBuffer)
	}
}

func (vp *vboPool) Borrow() (ret uint32) {
	vp.mu.Lock()
	defer vp.mu.Unlock()
	// move from one list to another
	for vp.vacantList.Len() == 0 {
		vp.mu.Unlock()
		// passive waiting graphic render loop to enlarge the queue
		time.Sleep(time.Millisecond)
		vp.mu.Lock()
	}

	allocate := vp.vacantList.Back()
	vp.vacantList.Remove(allocate)
	ret = allocate.Value.(uint32)
	vp.inuseIndex[ret] = vp.inuseList.PushBack(allocate.Value)
	return
}

func (vp *vboPool) Release(idx uint32) {
	vp.mu.Lock()
	node, ok := vp.inuseIndex[idx]
	if !ok {
		panic("cannot release an unexistent buffer uint32")
	}
	vp.inuseList.Remove(node)
	vp.vacantList.PushBack(node.Value)
	delete(vp.inuseIndex, idx)
	defer vp.mu.Unlock()

}

// Enlarge operation should be performed on Render thread,
// else it will cause panic when handling buffer allocation.
func (vp *vboPool) Enlarge(space int) {
	vp.mu.Lock()
	for i := 0; i < space; i++ {
		vp.vacantList.PushBack(GLNewVBO(1))
	}
	vp.mu.Unlock()
}

func (vp *vboPool) Len() (ret int) {
	ret = vp.vacantList.Len()
	return
}
