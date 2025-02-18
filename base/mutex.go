package base

import "sync/atomic"

// 尝试锁
type Mutex struct {
	state int32
}

// 尝试获取锁
func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapInt32(&m.state, 0, 1)
}

// 释放锁
func (m *Mutex) Unlock() {
	atomic.StoreInt32(&m.state, 0)
}
