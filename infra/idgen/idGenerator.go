package idgen

import (
	"fmt"

	"galaxyzeta.io/engine/infra/concurrency/lock"
)

type IdGenerator struct {
	Prefix string
	static int64
	mu     *lock.SpinLock
}

func NewIdGenerator(prefix string) *IdGenerator {
	return &IdGenerator{
		Prefix: prefix,
		static: 0,
	}
}

func (idgen *IdGenerator) Generate() string {
	idgen.mu.Lock()
	ret := fmt.Sprintf("%s%d", idgen.Prefix, idgen.static)
	idgen.static++
	idgen.mu.Unlock()
	return ret
}
