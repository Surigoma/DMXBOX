package fps_test

import (
	"backend/dmxServer/fps"
	"testing"
	"time"
)

func TestNewFPS(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		fps      float32
		callback func() bool
		finalize func()
		isNil    bool
	}{
		{
			name:     "Create",
			fps:      30,
			callback: func() bool { return true },
			finalize: func() {},
			isNil:    false,
		},
		{
			name:     "not callback",
			fps:      30,
			callback: nil,
			finalize: func() {},
			isNil:    true,
		},
		{
			name:     "not finalize",
			fps:      30,
			callback: func() bool { return true },
			finalize: nil,
			isNil:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := fps.NewFPS(tt.fps, tt.callback, tt.finalize) == nil
			t.Logf("result: %v", got)
			if got != tt.isNil {
				t.Errorf("NewFPS() = %v, want %v", got, tt.isNil)
			}
		})
	}
}

func TestFPSController_RunBreak(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		callback func() bool
		fps      float32
		want     bool
	}{
		{
			name:     "Return true",
			fps:      30,
			callback: func() bool { return true },
			want:     true,
		},
		{
			name:     "Return false",
			fps:      30,
			callback: func() bool { return false },
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			running := true
			fps := fps.NewFPS(tt.fps, tt.callback, func() { running = false })
			go fps.Run()
			time.Sleep(time.Duration(time.Second))
			t.Logf("result: %v", running)
			if running != tt.want {
				t.Errorf("Failed %v want %v", running, tt.want)
			}
		})
	}
}

func TestFPSController_Run(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		fps  float32
	}{
		{
			name: "10 FPS",
			fps:  10,
		},
		{
			name: "30 FPS",
			fps:  30,
		},
		{
			name: "50 FPS",
			fps:  50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var testChannel = make(chan bool)
			var callback = func() bool {
				testChannel <- true
				return true
			}
			var finalize = func() {
				testChannel <- true
			}
			fps := fps.NewFPS(tt.fps, callback, finalize)
			go fps.Run()
			select {
			case <-testChannel:
				break
			case <-time.After(time.Duration((1.1 / tt.fps) * float32(time.Second))):
				t.Error("Failed call callback function")
				return
			}
			fps.Stop()
			select {
			case <-testChannel:
				break
			case <-time.After(time.Duration(time.Second)):
				t.Error("Failed to call finalize")
				return
			}
		})
	}
}

func TestFPSController_GetFPS(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		fps  float32
	}{
		{
			name: "10 FPS",
			fps:  10,
		},
		{
			name: "30 FPS",
			fps:  30,
		},
		{
			name: "50 FPS",
			fps:  50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fps := fps.NewFPS(tt.fps, func() bool { return true }, func() {})
			go fps.Run()
			time.Sleep(time.Duration(time.Second))
			under := tt.fps * 0.95
			upper := tt.fps * 1.05
			got := fps.GetFPS()
			t.Logf("result: %v fps", got)
			if got < under && got > upper {
				t.Errorf("GetFPS() = %v, want %v", got, tt.fps)
			}
		})
	}
}
func TestFPSController_GetFPS_Notrun(t *testing.T) {
	fps := fps.NewFPS(30, func() bool { return true }, func() {})
	got := fps.GetFPS()
	if got > 0 {
		t.Errorf("Failed to test. %v want under 0", got)
	}
}
