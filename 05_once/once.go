package once

import (
	mutex "primitives/02_mutex"
	"primitives/internal/futex"
	"sync/atomic"
)

type Once struct {
	state uint32
	mutex mutex.Mutex
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.state) == 1 {
		return
	}
	o.mutex.Lock()
	defer o.mutex.Unlock() 

	if atomic.LoadUint32(&o.state) == 0 {
		defer atomic.StoreUint32(&o.state, 1)
		f()
	}
}

func (o *Once) Done() bool {
	if atomic.LoadUint32(&o.state) == 1{
		futex.WakeAll(&o.state)
		return true
	}
	return false
}
