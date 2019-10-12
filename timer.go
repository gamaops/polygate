package main

import (
	"sync/atomic"
	"time"
)

type ResetableTimerStatus uint8

const (
	RTCancel  ResetableTimerStatus = 0
	RTTimeout ResetableTimerStatus = 1
)

type ResetableTimer struct {
	resetCh     chan bool
	cancelCh    chan bool
	cancelCount int32
	duration    time.Duration
	Status      chan ResetableTimerStatus
}

func (r *ResetableTimer) Cancel() {
	if atomic.CompareAndSwapInt32(&r.cancelCount, 0, 1) {
		select {
		case r.cancelCh <- true:
		default:
		}
	}
}

func (r *ResetableTimer) Reset() {
	r.resetCh <- true
}

func (r *ResetableTimer) Start() {
	atomic.StoreInt32(&r.cancelCount, 0)
	select {
	case <-r.cancelCh:
		r.Status <- RTCancel
	case <-time.After(r.duration):
		r.Status <- RTTimeout
	case <-r.resetCh:
		go r.Start()
	}
}

func NewResetableTimer(duration time.Duration) *ResetableTimer {
	r := &ResetableTimer{
		duration:    duration,
		cancelCount: 0,
		resetCh:     make(chan bool, 1),
		cancelCh:    make(chan bool, 1),
		Status:      make(chan ResetableTimerStatus, 1),
	}

	go r.Start()

	return r
}
