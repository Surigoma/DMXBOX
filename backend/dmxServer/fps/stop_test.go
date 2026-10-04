package fps_test

import (
	"backend/dmxServer/fps"
	"sync/atomic"
	"testing"
	"time"
)

func TestStopBeforeRun(t *testing.T) {
	var calls atomic.Int32
	done := make(chan struct{})
	controller := fps.NewFPS(30, func() bool { calls.Add(1); return true }, func() { close(done) })
	controller.Stop()
	go controller.Run()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("stop was lost before Run")
	}
	if calls.Load() != 0 || controller.Running.Load() {
		t.Fatal("stopped controller ran a callback")
	}
}
