package mutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const (
	free = iota
	held
	contended
)

type Mutex struct {
	state uint32
}

func (m *Mutex) Lock() {
	for { 
		if atomic.CompareAndSwapUint32(&m.state, free, held) {
			return
		}
		for i := 0; i < 100; i++ {
			state := atomic.LoadUint32(&m.state)
			if state == free {
				if atomic.CompareAndSwapUint32(&m.state, free, held) {
					return 
				}
			}
		}

		atomic.CompareAndSwapUint32(&m.state, held, contended) 
		for {
			futex.Wait(&m.state, contended)
			if atomic.CompareAndSwapUint32(&m.state, free, contended) {
				return
			}
			atomic.CompareAndSwapUint32(&m.state, held, contended)
		}
	}
}

func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapUint32(&m.state, free, held)
}

func (m *Mutex) Unlock() {
	oldState := atomic.SwapUint32(&m.state, free)
	if oldState == held {
		return
	}
	if oldState == contended {
		futex.Wake(&m.state)
		return
	}
	panic("A")
}
