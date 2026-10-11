package waitgroup

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type WaitGroup struct {
	count uint32
}

func (wg *WaitGroup) Add(delta int) {
	for {
		old_v := atomic.LoadUint32(&wg.count)
		if delta < 0 && int64(old_v)+int64(delta) < 0 {
			panic("negative counter")
		}
		new_v := old_v + uint32(delta)
		if !atomic.CompareAndSwapUint32(&wg.count, old_v, new_v) {
			continue
		}
		if new_v == 0 {
			futex.WakeAll(&wg.count)
		}
		return
	}
}

func (wg *WaitGroup) Done() {
	new_v := atomic.AddUint32(&wg.count, ^uint32(0))
	if new_v == 0 {
		futex.WakeAll(&wg.count)
	}
	if new_v == ^uint32(0) {
		panic("Bad value")
	}
}

func (wg *WaitGroup) Wait() {
	for {
		v := atomic.LoadUint32(&wg.count)
		if v == 0 {
			return
		}
		futex.Wait(&wg.count, v)
	}
}
