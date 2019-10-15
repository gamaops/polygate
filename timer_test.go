package main

import (
	"testing"
	"time"
)

func TestNewRestableTimer(t *testing.T) {

	r := NewResetableTimer(1 * time.Second)

	if r == nil {
		t.Error("NewResetableTimer must return a pointer to ResetableTimer")
	}

}

func TestResetableTimerOperations(t *testing.T) {

	r := NewResetableTimer(100 * time.Millisecond)

	if r.cancelCount != 0 {
		t.Error("Cancel count of ResetableTimer must be set to zero")
	}

	r.Cancel()

	st := <-r.Status

	if st != RTCancel {
		t.Error("Invalid status for cancel operation")
	}

	go r.Start()

	time.Sleep(120 * time.Millisecond)

	st = <-r.Status

	if st != RTTimeout {
		t.Error("ResetableTimer must time out after duration")
	}

	go r.Start()

	time.Sleep(90 * time.Millisecond)

	r.Reset()

	time.Sleep(90 * time.Millisecond)

	r.Cancel()

	st = <-r.Status

	if st != RTCancel {
		t.Error("ResetableTimer reset the timer calling Reset")
	}

}
