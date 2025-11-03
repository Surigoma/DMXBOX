package module_test

import (
	"backend/config"
	"backend/dmxServer/controller"
	"backend/dmxServer/controller/module"
	"log/slog"
	"runtime"
	"sync"
	"testing"

	"go.bug.st/serial"
)

var testData = map[string]func() *controller.Controller{
	"artnet":  module.NewArtnet,
	"ftdi":    module.NewFTDI,
	"console": module.NewConsole,
}

func TestModule_New(t *testing.T) {
	tests := []string{}
	for k := range testData {
		tests = append(tests, k)
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			got := testData[tt]()
			if got.ModInitialize == nil {
				t.Error("ModInitialize is nil")
				return
			}
			if got.ModOutput == nil {
				t.Error("ModOutput is nil")
				return
			}
			if got.ModFinalize == nil {
				t.Error("ModFinalize is nil")
				return
			}
			t.Logf("model: %v", got)
		})
	}
}

func getConsolePorts() []string {
	consoles, err := serial.GetPortsList()
	if err != nil {
		slog.Error("Can not get console port.", "err", err)
		return []string{}
	}
	canOpenPort := []string{}
	for _, name := range consoles {
		port, err := serial.Open(name, &serial.Mode{
			BaudRate: 250000,
			DataBits: 8,
			Parity:   serial.NoParity,
		})
		if err == nil {
			port.Close()
			canOpenPort = append(canOpenPort, name)
		}
	}
	if len(canOpenPort) <= 0 {
		slog.Error("Not found for can open console.")
		GOOS := runtime.GOOS
		switch GOOS {
		case "windows":
			slog.Error("Can use com0com (https://sourceforge.net/projects/com0com/)")
		case "linux":
			slog.Error("Can use socat (https://linux.die.net/man/1/socat)")
		}
		return []string{}
	}
	return canOpenPort
}

func TestModule_Initialize(t *testing.T) {
	canOpenPort := getConsolePorts()
	if len(canOpenPort) <= 0 {
		t.Error("Failed to get ports.")
		return
	}
	tests := []struct {
		name   string
		target string
		config config.Config
		want   bool
	}{
		{
			name:   "Console",
			target: "console",
			config: config.Config{
				Output: config.OutputTargets{
					Target: []string{"console"},
				},
			},
			want: true,
		},
		{
			name:   "FTDI",
			target: "ftdi",
			config: config.Config{
				Output: config.OutputTargets{
					Target: []string{"ftdi"},
					DMX: config.DMXHardware{
						Port: canOpenPort[0],
					},
				},
			},
			want: true,
		},
		{
			name:   "FTDI Undefined port",
			target: "ftdi",
			config: config.Config{
				Output: config.OutputTargets{
					Target: []string{"ftdi"},
					DMX: config.DMXHardware{
						Port: "COM10000",
					},
				},
			},
			want: false,
		},
		{
			name:   "Artnet",
			target: "artnet",
			config: config.Config{
				Output: config.OutputTargets{
					Target: []string{"artnet"},
					Artnet: config.Artnet{
						Address:     "127.0.0.1",
						Universe:    0,
						SubUniverse: 0,
						Net:         0,
					},
				},
			},
			want: true,
		},
		{
			name:   "Artnet Unsupport address",
			target: "artnet",
			config: config.Config{
				Output: config.OutputTargets{
					Target: []string{"artnet"},
					Artnet: config.Artnet{
						Address:     "254.254.254.254",
						Universe:    0,
						SubUniverse: 0,
						Net:         0,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
			cont := testData[tt.target]()
			got := cont.Initialize(&tt.config, logger)
			defer cont.Finalize()
			if got != tt.want {
				t.Errorf("Failed to initialize: %v", cont)
			}
		})
	}
}

func TestModule_FTDI_Output(t *testing.T) {
	var lock sync.Mutex
	canOpenPort := getConsolePorts()
	if len(canOpenPort) <= 0 {
		t.Error("Failed to get ports.")
		return
	}
	t.Run("Can Output", func(tt *testing.T) {
		lock.Lock()
		defer lock.Unlock()
		logger := slog.New(slog.NewJSONHandler(tt.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
		target := canOpenPort[0]
		ftdi := testData["ftdi"]()
		config := config.Config{
			Output: config.OutputTargets{
				Target: []string{"ftdi"},
				DMX: config.DMXHardware{
					Port: target,
				},
			},
		}
		if !ftdi.Initialize(&config, logger) {
			tt.Error("Failed to initialize.")
			return
		}
		defer ftdi.Finalize()
		test := make([]byte, 512)
		if !ftdi.Output(&test) {
			tt.Error("Failed to output.")
			return
		}
	})
	t.Run("Port is down", func(tt *testing.T) {
		lock.Lock()
		defer lock.Unlock()
		logger := slog.New(slog.NewJSONHandler(tt.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
		target := canOpenPort[0]
		ftdi := testData["ftdi"]()
		config := config.Config{
			Output: config.OutputTargets{
				Target: []string{"ftdi"},
				DMX: config.DMXHardware{
					Port: target,
				},
			},
		}
		if !ftdi.Initialize(&config, logger) {
			tt.Error("Failed to initialize.")
			return
		}
		ftdi.Finalize()
		test := make([]byte, 512)
		if ftdi.Output(&test) { // Port error
			tt.Error("Failed to close output.")
			return
		}
		if ftdi.Output(&test) { // Nil
			tt.Error("Failed to close output.")
			return
		}
	})
}

func TestModule_Console_Output(t *testing.T) {
	t.Run("Console", func(tt *testing.T) {
		logger := slog.New(slog.NewJSONHandler(tt.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
		console := testData["console"]()
		config := config.Config{
			Output: config.OutputTargets{
				Target: []string{"console"},
			},
		}
		if !console.Initialize(&config, logger) {
			tt.Error("Failed to initialize.")
			return
		}
		defer console.Finalize()
		test := make([]byte, 512)
		if !console.Output(&test) {
			tt.Error("Failed to output.")
			return
		}
	})
}

func TestModule_Artnet_Output(t *testing.T) {
	t.Run("Artnet", func(tt *testing.T) {
		logger := slog.New(slog.NewJSONHandler(tt.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
		artnet := testData["artnet"]()
		config := config.Config{
			Output: config.OutputTargets{
				Target: []string{"artnet"},
				Artnet: config.Artnet{
					Address:     "127.0.0.1",
					Universe:    0,
					SubUniverse: 0,
					Net:         0,
				},
			},
		}
		if !artnet.Initialize(&config, logger) {
			tt.Error("Failed to initialize.")
			return
		}
		defer artnet.Finalize()
		test := make([]byte, 512)
		if !artnet.Output(&test) {
			tt.Error("Failed to output.")
			return
		}
	})
}
