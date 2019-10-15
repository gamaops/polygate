package main

import (
	"time"
)

type SafePool struct {
	transfer   chan interface{}
	New        func() (interface{}, error)
	Invalidate func(interface{})
}

func NewSafePool() *SafePool {
	return &SafePool{
		transfer: make(chan interface{}),
	}
}

func (s *SafePool) Get(timeout time.Duration) (interface{}, error) {
	select {
	case item := <-s.transfer:
		return item, nil
	case <-time.After(timeout):
		return s.New()
	}
}

func (s *SafePool) Put(item interface{}, timeout time.Duration) {

	select {
	case s.transfer <- item:
		return
	case <-time.After(timeout):
		s.Invalidate(item)
	}

}
