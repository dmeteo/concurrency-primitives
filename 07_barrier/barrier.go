package barrier

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Barrier struct {
	need    uint32
	arrived uint32
	round   uint32
}

func New(n int) *Barrier {
	return &Barrier{need: uint32(n), arrived: 0, round: 0}
}

func (b *Barrier) Wait() {
	round := atomic.LoadUint32(&b.round)
	need := atomic.LoadUint32(&b.need)
	arr := atomic.AddUint32(&b.arrived, uint32(1))
	if need > arr {
		for {
			if atomic.LoadUint32(&b.round) != round {
				return
			}
			futex.Wait(&b.round, round)	
		}
	}
	if need == arr {
		atomic.SwapUint32(&b.arrived, 0)
		atomic.AddUint32(&b.round, uint32(1))
		futex.WakeAll(&b.round)
	}
}
