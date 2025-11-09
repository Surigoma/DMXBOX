package dmxserver_test

import (
	"backend/config"
	dmxserver "backend/dmxServer"
	"backend/dmxServer/controller"
	device "backend/dmxServer/devices"
	"backend/packageModule"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"testing"
	"time"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name   string
		module *packageModule.PackageModule
		config *config.Config
		want   bool
	}{
		{
			name:   "Success",
			module: &packageModule.PackageModule{},
			config: &config.Config{
				Output: config.OutputTargets{
					Target: []string{"console"},
				},
				Dmx: config.DMXServer{
					Groups: map[string]config.DMXGroup{
						"test": {
							Name: "TEST",
							Devices: []config.DMXDevice{
								{
									Model:    "dimmer",
									Channel:  1,
									MaxValue: []uint{255},
								},
								{
									Model:    "wclight",
									Channel:  2,
									MaxValue: []uint{255, 255, 255},
								},
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name:   "Success - duplicate target",
			module: &packageModule.PackageModule{},
			config: &config.Config{
				Output: config.OutputTargets{
					Target: []string{"console", "console"},
				},
				Dmx: config.DMXServer{
					Groups: map[string]config.DMXGroup{
						"test": {
							Name: "TEST",
							Devices: []config.DMXDevice{
								{
									Model:    "dimmer",
									Channel:  1,
									MaxValue: []uint{255},
								},
								{
									Model:    "wclight",
									Channel:  2,
									MaxValue: []uint{255, 255, 255},
								},
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name:   "Success - gruop is 0",
			module: &packageModule.PackageModule{},
			config: &config.Config{
				Output: config.OutputTargets{
					Target: []string{"console"},
				},
				Dmx: config.DMXServer{
					Groups: map[string]config.DMXGroup{},
				},
			},
			want: true,
		},
		{
			name:   "Success - device is 0",
			module: &packageModule.PackageModule{},
			config: &config.Config{
				Output: config.OutputTargets{
					Target: []string{"console"},
				},
				Dmx: config.DMXServer{
					Groups: map[string]config.DMXGroup{
						"test": {
							Name:    "TEST",
							Devices: []config.DMXDevice{},
						},
					},
				},
			},
			want: true,
		},
		{
			name:   "Failed - Unknown model",
			module: &packageModule.PackageModule{},
			config: &config.Config{
				Output: config.OutputTargets{
					Target: []string{"console"},
				},
				Dmx: config.DMXServer{
					Groups: map[string]config.DMXGroup{
						"test": {
							Name: "TEST",
							Devices: []config.DMXDevice{
								{
									Model:    "UNKNOWN_MODEL",
									Channel:  1,
									MaxValue: []uint{255},
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name:   "Failed - Unknown Output target",
			module: &packageModule.PackageModule{},
			config: &config.Config{
				Output: config.OutputTargets{
					Target: []string{"UNKNOWN"},
				},
				Dmx: config.DMXServer{
					Groups: map[string]config.DMXGroup{
						"test": {
							Name: "TEST",
							Devices: []config.DMXDevice{
								{
									Model:    "dimmer",
									Channel:  1,
									MaxValue: []uint{255},
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name:   "Failed - controller is nil",
			module: &packageModule.PackageModule{},
			config: &config.Config{
				Output: config.OutputTargets{
					Target: []string{},
				},
				Dmx: config.DMXServer{
					Groups: map[string]config.DMXGroup{
						"test": {
							Name: "TEST",
							Devices: []config.DMXDevice{
								{
									Model:    "dimmer",
									Channel:  1,
									MaxValue: []uint{255},
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name:   "Failed - mismatch max value for device",
			module: &packageModule.PackageModule{},
			config: &config.Config{
				Output: config.OutputTargets{
					Target: []string{"console"},
				},
				Dmx: config.DMXServer{
					Groups: map[string]config.DMXGroup{
						"test": {
							Name: "TEST",
							Devices: []config.DMXDevice{
								{
									Model:    "dimmer",
									Channel:  1,
									MaxValue: []uint{},
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name:   "Failed - failed to setup controller",
			module: &packageModule.PackageModule{},
			config: &config.Config{
				Output: config.OutputTargets{
					Target: []string{"artnet"},
					Artnet: config.Artnet{
						Address:     "254.254.254.254",
						Universe:    0,
						SubUniverse: 0,
						Net:         0,
					},
				},
				Dmx: config.DMXServer{
					Groups: map[string]config.DMXGroup{
						"test": {
							Name: "TEST",
							Devices: []config.DMXDevice{
								{
									Model:    "dimmer",
									Channel:  1,
									MaxValue: []uint{},
								},
							},
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.module.Logger = slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
			o, _ := json.Marshal(tt.config)
			t.Logf("config %v", string(o))
			got := dmxserver.DMXServer.Initialize(tt.module, tt.config)
			if got != tt.want {
				t.Errorf("Initialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareConfig(v1 config.DMXGroup, v2 config.DMXGroup) bool {
	if v1.Name != v2.Name {
		return false
	}
	if !slices.EqualFunc(v1.Devices, v2.Devices, func(v1 config.DMXDevice, v2 config.DMXDevice) bool {
		if v1.Model != v2.Model {
			return false
		}
		if v1.Channel != v2.Channel {
			return false
		}
		if !slices.Equal(v1.MaxValue, v2.MaxValue) {
			return false
		}
		return true
	}) {
		return false
	}
	return true
}

func TestGetConfig(t *testing.T) {
	group := map[string]config.DMXGroup{}
	index := 1
	for gi := range 10 {
		key := fmt.Sprintf("test%d", gi)
		name := fmt.Sprintf("TEST_%d", gi)
		group[key] = config.DMXGroup{
			Name:    name,
			Devices: make([]config.DMXDevice, 20),
		}
		for di := range 20 {
			group[key].Devices[di] = config.DMXDevice{
				Model:    "dimmer",
				Channel:  uint8(index),
				MaxValue: []uint{255},
			}
			index++
		}
	}
	tests := []struct {
		name   string
		config config.Config
	}{
		{
			name: "Conpare",
			config: config.Config{
				Output: config.OutputTargets{
					Target: []string{"console"},
				},
				Dmx: config.DMXServer{
					Groups: group,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.config.Dmx.Groups
			module := packageModule.PackageModule{
				Logger: slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelInfo})),
			}
			if !dmxserver.Initialize(&module, &tt.config) {
				t.Error("Failed to initialize dmx server.")
				return
			}
			got := dmxserver.GetConfig()
			gotJson, _ := json.Marshal(got)
			if !maps.EqualFunc(got, want, compareConfig) {
				wantJson, _ := json.Marshal(want)
				t.Errorf("GetConfig() = %v, want %v", string(gotJson), string(wantJson))
			}
			t.Logf("Success: %s", string(gotJson))
		})
	}
}

func TestRender(t *testing.T) {
	testConfig := config.Config{
		Output: config.OutputTargets{
			Target: []string{"test"},
		},
		Dmx: config.DMXServer{
			Groups: map[string]config.DMXGroup{
				"test": {
					Name: "test",
					Devices: []config.DMXDevice{
						{
							Model:    "test",
							Channel:  1,
							MaxValue: []uint{255},
						},
					},
				},
			},
		},
	}
	dmxserver.RenderTypes["test"] = func() *controller.Controller {
		return &controller.Controller{
			Model:         "test",
			ModInitialize: func(c *config.Config, l *slog.Logger) bool { return true },
			ModOutput:     func(b *[]byte) bool { return true },
			ModFinalize:   func() {},
		}
	}
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "device return true",
			want: true,
		},
		{
			name: "device return false",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
			module := packageModule.PackageModule{
				Logger: logger,
			}
			testChan := make(chan bool)
			dmxserver.DeviceTypes["test"] = func() *device.DMXDevice {
				return &device.DMXDevice{
					Model:      "test",
					UseChannel: 1,
					ModUpdate: func() bool {
						testChan <- tt.want
						return tt.want
					},
				}
			}
			if !dmxserver.Initialize(&module, &testConfig) {
				t.Error("failed to initialize dmx server")
				return
			}
			got := false
			go func() {
				got = dmxserver.Render()
			}()
			select {
			case <-testChan:
				break
			case <-time.After(time.Duration(time.Second)):
				t.Error("Failed to catch channel.")
				return
			}

			if got != tt.want {
				t.Errorf("Render() = %v, want %v", got, tt.want)
			}
		})
	}
}
