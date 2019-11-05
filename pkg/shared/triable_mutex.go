package shared

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const mutexLocked = 1 << iota

type TriableMutex struct {
	sync.Mutex
}

func (m *TriableMutex) TryLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked)
}
