package main

import (
	"testing"
	"time"
)

func TestNewSafePool(t *testing.T) {

	r := NewSafePool()

	if r == nil {
		t.Error("NewSafePool must return a pointer to SafePool")
	}

}

func TestSafePoolOperations(t *testing.T) {

	newValue := 1
	r := NewSafePool()

	var invalidateCalled *int = nil

	r.New = func() (interface{}, error) {
		return &newValue, nil
	}

	r.Invalidate = func(val interface{}) {
		invalidateCalled = val.(*int)
	}

	x, err := r.Get(10 * time.Millisecond)
	y := x.(*int)

	if err != nil {
		t.Error("SafePool must no return error")
	}
	if y != &newValue {
		t.Error("SafePool must return valid value (ptr)")
	}
	if *y != newValue {
		t.Error("SafePool must return valid value")
	}
	if invalidateCalled != nil {
		t.Error("SafePool must not invalidate valid items")
	}

	go r.Put(&newValue, 10*time.Millisecond)
	r.New = nil
	newValue++
	x, err = r.Get(10 * time.Millisecond)
	y = x.(*int)

	if err != nil {
		t.Error("SafePool must no return error")
	}
	if y != &newValue {
		t.Error("SafePool must return valid value (ptr)")
	}
	if *y != newValue {
		t.Error("SafePool must return valid value")
	}
	if invalidateCalled != nil {
		t.Error("SafePool must not invalidate valid items")
	}

	r.Put(&newValue, 10*time.Millisecond)

	if invalidateCalled != &newValue {
		t.Error("SafePool must return valid value (ptr)")
	}
	if *invalidateCalled != newValue {
		t.Error("SafePool must return valid value")
	}

}
