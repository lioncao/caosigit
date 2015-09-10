package tools

import (
	"runtime"
	"sync/atomic"
)

const (
	UNLOCKED = 0
	LOCKED   = 1
)

const (
	MAX_LOOP_COUNT = 100000
)

type FastLock struct {
	flag      int32
	loopCount int32
}

func NewFastLock() *FastLock {
	this := new(FastLock)
	this.flag = UNLOCKED
	this.loopCount = MAX_LOOP_COUNT
	return this
}

func NewFastLockWithLoopCount(loopCount int32) *FastLock {
	this := new(FastLock)
	this.flag = UNLOCKED
	if loopCount <= 0 {
		loopCount = MAX_LOOP_COUNT
	}
	this.loopCount = loopCount
	return this
}

func (this *FastLock) Lock() bool {

	for i := int32(0); i < this.loopCount; i++ {
		if atomic.CompareAndSwapInt32(&this.flag, UNLOCKED, LOCKED) {
			return true
		}
		runtime.Gosched()
	}
	return false
}

func (this *FastLock) Unlock() bool {
	for i := int32(0); i < this.loopCount; i++ {
		if atomic.CompareAndSwapInt32(&this.flag, LOCKED, UNLOCKED) {
			return true
		}
		runtime.Gosched()
	}
	return false
}
