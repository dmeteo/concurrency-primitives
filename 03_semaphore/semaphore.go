package semaphore

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Semaphore struct {
	permits uint32
}

func New(n int) *Semaphore {
	return &Semaphore{permits: uint32(n)}
}

func (s *Semaphore) Acquire() {
	for {
		if s.TryAcquire() {
			return
		}
		futex.Wait(&s.permits, 0)
	}
}

func (s *Semaphore) TryAcquire() bool {
	old := atomic.LoadUint32(&s.permits)
	if old == 0 {
		return false
	}
	new := old - 1
	return atomic.CompareAndSwapUint32(&s.permits, old, new)
}

func (s *Semaphore) Release() {
	for {
		old := atomic.LoadUint32(&s.permits)
		new := old + uint32(1)
		if atomic.CompareAndSwapUint32(&s.permits, old, new){ 
			futex.Wake(&s.permits)
			return
		}
	}
}

func (s *Semaphore) Available() int {
	return int(atomic.LoadUint32(&s.permits))
}
