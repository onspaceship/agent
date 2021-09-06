package socket

import "sync/atomic"

type refCounter struct {
	ref *int64
}

func newRefCounter() *refCounter {
	return &refCounter{ref: new(int64)}
}

func (counter *refCounter) nextRef() int64 {
	return atomic.AddInt64(counter.ref, 1)
}
