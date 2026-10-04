package config_test

import (
	"backend/config"
	"math"
	"os"
	"reflect"
	"testing"
)

func TestValidateUnsafeSettings(t *testing.T) {
	tests := []struct {
		name   string
		change func(*config.Config)
	}{
		{"zero FPS", func(c *config.Config) { c.Dmx.Fps = 0 }},
		{"NaN FPS", func(c *config.Config) { c.Dmx.Fps = float32(math.NaN()) }},
		{"negative delay", func(c *config.Config) { c.Dmx.Delay = -1 }},
		{"unknown input", func(c *config.Config) { c.Input.Modules = []string{"invalid"} }},
		{"unknown output", func(c *config.Config) { c.Output.Target = []string{"invalid"} }},
		{"unknown OSC type", func(c *config.Config) { c.Output.Target = []string{"console", "osc"}; c.Output.Osc.Type = "invalid" }},
		{"zero channel", func(c *config.Config) {
			c.Dmx.Groups["test"] = config.DMXGroup{Devices: []config.DMXDevice{{Model: "dimmer", Channel: 0, MaxValue: []uint{255}}}}
		}},
		{"overflow channels", func(c *config.Config) {
			c.Dmx.Groups["test"] = config.DMXGroup{Devices: []config.DMXDevice{{Model: "wclight", Channel: 511, MaxValue: []uint{255, 255, 0}}}}
		}},
		{"overflow max", func(c *config.Config) {
			c.Dmx.Groups["test"] = config.DMXGroup{Devices: []config.DMXDevice{{Model: "dimmer", Channel: 1, MaxValue: []uint{256}}}}
		}},
		{"short max", func(c *config.Config) {
			c.Dmx.Groups["test"] = config.DMXGroup{Devices: []config.DMXDevice{{Model: "wclight", Channel: 1, MaxValue: []uint{255}}}}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.InitializeConfig()
			c := config.Get()
			tt.change(&c)
			if c.Validate() == nil {
				t.Fatal("unsafe configuration was accepted")
			}
		})
	}
	config.InitializeConfig()
	c := config.Get()
	c.Dmx.Groups["edge"] = config.DMXGroup{Devices: []config.DMXDevice{{Model: "wclight", Channel: 510, MaxValue: []uint{255, 255, 0}}}}
	if err := c.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestSaveAndSetFailurePreservesRunningSettings(t *testing.T) {
	t.Chdir(t.TempDir())
	config.InitializeConfig()
	before := config.Get()
	changed := before
	changed.Dmx.Fps = 60
	if err := os.Mkdir("config.json", 0755); err != nil {
		t.Fatal(err)
	}
	if err := config.SaveAndSet(changed); err == nil {
		t.Fatal("save should fail")
	}
	if !reflect.DeepEqual(before, config.Get()) {
		t.Fatal("failed save changed running configuration")
	}
}
