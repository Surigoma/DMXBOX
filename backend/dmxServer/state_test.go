package dmxserver

import (
	"backend/config"
	"backend/message"
	"backend/packageModule"
	"io"
	"log/slog"
	"math"
	"sync"
	"testing"
)

func TestFadeStates(t *testing.T) {
	config.InitializeConfig()
	c := config.Get()
	c.Dmx.Groups = map[string]config.DMXGroup{
		"stage": {Devices: []config.DMXDevice{
			{Model: "dimmer", Channel: 1, MaxValue: []uint{255}},
			{Model: "wclight", Channel: 5, MaxValue: []uint{128, 255, 0}},
		}},
		"empty": {},
		"zero":  {Devices: []config.DMXDevice{{Model: "dimmer", Channel: 11, MaxValue: []uint{0}}}},
	}
	module := packageModule.PackageModule{Logger: slog.New(slog.NewTextHandler(io.Discard, nil)), Wg: &sync.WaitGroup{}}
	if !Initialize(&module, &c) {
		t.Fatal("initialize failed")
	}
	defer CleanupDMXServer()
	fade := func(isIn, duration, interval string) {
		handleMessage(message.Message{Arg: message.MessageBody{Action: "fade", Arg: map[string]string{"id": "stage", "isIn": isIn, "duration": duration, "interval": interval}}})
	}
	if got := GetFadeStates()["stage"]; got.Level != 0 || got.State != "idle" {
		t.Fatal(got)
	}
	fade("true", "1", "60")
	if got := GetFadeStates()["stage"]; got.Level != 0 || got.State != "waiting" || !got.IsIn {
		t.Fatal(got)
	}
	fade("true", "1", "0")
	// Read the actual output, rather than predicting a level from elapsed time.
	stateMutex.Lock()
	rendered[0], rendered[4], rendered[5], rendered[6] = 128, 64, 128, 255
	stateMutex.Unlock()
	if got := GetFadeStates()["stage"]; math.Abs(got.Level-320.0/638.0) > 0.0001 || got.State != "fading" || !got.IsIn {
		t.Fatal(got)
	}
	fade("false", "1", "0")
	if got := GetFadeStates()["stage"]; got.State != "fading" || got.IsIn || got.Level == 0 {
		t.Fatal(got)
	}
	fade("true", "0", "0")
	Render()
	if got := GetFadeStates()["stage"]; got.Level != 1 || got.State != "idle" {
		t.Fatal(got)
	}
	fade("false", "0", "0")
	Render()
	if got := GetFadeStates()["stage"]; got.Level != 0 || got.State != "idle" {
		t.Fatal(got)
	}
	for _, id := range []string{"empty", "zero"} {
		if got := GetFadeStates()[id]; got.Level != 0 || got.State != "idle" {
			t.Fatal(id, got)
		}
	}
}
