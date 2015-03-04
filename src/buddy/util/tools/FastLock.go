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
	LOOP_COUNT = 100000
)

type FastLock struct {
	flag int32
}

func NewFastLock() *FastLock {
	this := new(FastLock)
	this.flag = UNLOCKED
	return this
}

func (this *FastLock) Lock() bool {

	for i := 0; i < LOOP_COUNT; i++ {
		if atomic.CompareAndSwapInt32(&this.flag, UNLOCKED, LOCKED) {
			return true
		}
		runtime.Gosched()
	}
	return false
}

func (this *FastLock) Unlock() bool {
	for i := 0; i < LOOP_COUNT; i++ {
		if atomic.CompareAndSwapInt32(&this.flag, LOCKED, UNLOCKED) {
			return true
		}
		runtime.Gosched()
	}
	return false
}
