package device_test

import (
	device "backend/dmxServer/devices"
	"bytes"
	"sync"
	"testing"
	"time"
)

func TestDMXDevice_Initialize(t *testing.T) {
	var duration float32 = 0.1
	tests := []struct {
		name     string
		channel  uint8
		maxValue []byte
		duration *float32
		want     bool
	}{
		{
			name:     "Can create",
			channel:  1,
			maxValue: []byte{255, 255, 255},
			duration: &duration,
			want:     true,
		},
		{
			name:     "Can not create",
			channel:  1,
			maxValue: []byte{},
			duration: &duration,
			want:     false,
		},
		{
			name:     "maxValue is nil",
			channel:  1,
			maxValue: nil,
			duration: &duration,
			want:     false,
		},
	}
	target := make([]byte, 512)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dev = device.DMXDevice{
				Model:      "test",
				UseChannel: 3,
			}
			got := dev.Initialize(tt.channel, tt.maxValue, &target, tt.duration)
			if got != tt.want {
				t.Errorf("Initialize() = %v, want %v", got, tt.want)
				return
			}
			t.Logf("Success %v wants %v", got, tt.want)
		})
	}
}

func TestDMXDevice_Fade(t *testing.T) {
	tests := []struct {
		name        string
		channel     uint8
		maxValue    []byte
		duration    float32
		addTime     float32
		optDuration float32
		optInterval float32
		isIn        bool
		modFade     bool
	}{
		{
			name:        "Fade In",
			channel:     1,
			maxValue:    []byte{255, 255, 255},
			duration:    0.8,
			addTime:     0,
			optDuration: -1,
			optInterval: -1,
			isIn:        true,
		},
		{
			name:        "Fade Out",
			channel:     1,
			maxValue:    []byte{255, 255, 255},
			duration:    0.8,
			addTime:     0,
			optDuration: -1,
			optInterval: -1,
			isIn:        false,
		},
		{
			name:        "Fade In (to target)",
			channel:     1,
			maxValue:    []byte{100, 200, 255},
			duration:    0.8,
			addTime:     0,
			optDuration: -1,
			optInterval: -1,
			isIn:        true,
		},
		{
			name:        "Fade In (Long time 2s)",
			channel:     1,
			maxValue:    []byte{100, 200, 255},
			duration:    2,
			addTime:     0,
			optDuration: -1,
			optInterval: -1,
			isIn:        true,
		},
		{
			name:        "Fade In (Shot time 0.01s)",
			channel:     1,
			maxValue:    []byte{100, 200, 255},
			duration:    0.01,
			addTime:     0,
			optDuration: -1,
			optInterval: -1,
			isIn:        true,
		},
		{
			name:        "Fade In (Overtime)",
			channel:     1,
			maxValue:    []byte{100, 200, 255},
			duration:    0.1,
			addTime:     0.5,
			optDuration: -1,
			optInterval: -1,
			isIn:        true,
		},
		{
			name:        "Fade In (Optional Duration 0.1s)",
			channel:     1,
			maxValue:    []byte{100, 200, 255},
			duration:    0.5,
			addTime:     0.5,
			optDuration: 0.1,
			optInterval: -1,
			isIn:        true,
		},
		{
			name:        "Fade In (Optional Duration 0s)",
			channel:     1,
			maxValue:    []byte{100, 200, 255},
			duration:    0.5,
			addTime:     0.5,
			optDuration: 0,
			optInterval: -1,
			isIn:        true,
		},
		{
			name:        "Fade In (Optional Interval 0.5s)",
			channel:     1,
			maxValue:    []byte{100, 200, 255},
			duration:    0.1,
			addTime:     0,
			optDuration: -1,
			optInterval: 0.5,
			isIn:        true,
		},
		{
			name:        "Change channel",
			channel:     10,
			maxValue:    []byte{100, 200, 255},
			duration:    0,
			addTime:     0,
			optDuration: -1,
			optInterval: -1,
			isIn:        true,
		},
		{
			name:        "ModFade - Fill max",
			channel:     1,
			maxValue:    []byte{100, 200, 255},
			duration:    0,
			addTime:     0,
			optDuration: -1,
			optInterval: -1,
			isIn:        true,
			modFade:     true,
		},
		{
			name:        "ModFade - Fill 0",
			channel:     1,
			maxValue:    []byte{100, 200, 255},
			duration:    0,
			addTime:     0,
			optDuration: -1,
			optInterval: -1,
			isIn:        false,
			modFade:     false,
		},
	}
	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			t.Parallel()
			target := make([]byte, 512)
			var dev device.DMXDevice = device.DMXDevice{
				Model:      "test",
				UseChannel: 3,
			}
			if !dev.Initialize(tt.channel, tt.maxValue, &target, &tt.duration) {
				t.Error("Failed to initialize")
			}
			if tt.modFade {
				dev.ModFade = func(isIn bool, duration float32, interval float32) {
					for i := range dev.Target {
						if isIn {
							(*dev.Output)[i+int(dev.Channel)-1] = dev.MaxValue[i]
						} else {
							(*dev.Output)[i+int(dev.Channel)-1] = 0
						}
					}
				}
			}
			dev.Fade(tt.isIn, tt.optDuration, tt.optInterval)
			if tt.optInterval > 0 {
				wg.Add(1)
				dev.Update(&wg)
				if !bytes.Equal(target[dev.Channel-1:dev.Channel+dev.UseChannel-1], make([]byte, dev.UseChannel)) {
					t.Error("Failed to match 0 when before interval")
					return
				}
				time.Sleep(time.Duration(tt.optInterval * float32(time.Second)))
			}
			if tt.optDuration < 0 {
				time.Sleep(time.Duration(tt.duration * float32(time.Second)))
			} else {
				time.Sleep(time.Duration(tt.optDuration * float32(time.Second)))
			}
			wg.Add(1)
			dev.Update(&wg)
			if tt.addTime > 0 {
				t.Log("Wait after")
				time.Sleep(time.Duration(tt.addTime * float32(time.Second)))
				wg.Add(1)
				dev.Update(&wg)
			}
			targetValue := tt.maxValue
			if !tt.isIn {
				targetValue = make([]byte, len(tt.maxValue))
			}
			if !bytes.Equal(target[tt.channel-1:tt.channel+dev.UseChannel-1], targetValue) {
				t.Errorf("Failed to fade. %v want %v", target[tt.channel-1:tt.channel+dev.UseChannel-1], targetValue)
				return
			}
			t.Logf("result: %v want %v", target[tt.channel-1:tt.channel+dev.UseChannel-1], targetValue)
		})
	}
}

func TestDMXDevice_Update(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		wg   *sync.WaitGroup
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var dev device.DMXDevice
			got := dev.Update(tt.wg)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
