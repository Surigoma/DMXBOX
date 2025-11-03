package controller_test

import (
	"backend/config"
	"backend/dmxServer/controller"
	"log/slog"
	"testing"
)

var testTarget = map[string]controller.Controller{
	"full": {
		Model:         "full",
		ModInitialize: func(c *config.Config, l *slog.Logger) bool { return true },
		ModOutput:     func(b *[]byte) bool { return true },
		ModFinalize:   func() {},
	},
	"non_init": {
		Model:         "non_init",
		ModInitialize: nil,
		ModOutput:     func(b *[]byte) bool { return true },
		ModFinalize:   func() {},
	},
	"non_output": {
		Model:         "non_output",
		ModInitialize: func(c *config.Config, l *slog.Logger) bool { return true },
		ModOutput:     nil,
		ModFinalize:   func() {},
	},
	"non_final": {
		Model:         "non_output",
		ModInitialize: func(c *config.Config, l *slog.Logger) bool { return true },
		ModOutput:     func(b *[]byte) bool { return true },
		ModFinalize:   nil,
	},
}

func TestController_Initialize(t *testing.T) {
	tests := []struct {
		name   string
		target controller.Controller
		config *config.Config
		want   bool
	}{
		{
			name:   "Full function",
			target: testTarget["full"],
			config: nil,
			want:   true,
		},
		{
			name:   "Non init function",
			target: testTarget["non_init"],
			config: nil,
			want:   false,
		},
		{
			name:   "Non output function",
			target: testTarget["non_output"],
			config: nil,
			want:   true,
		},
		{
			name:   "Non final function",
			target: testTarget["non_final"],
			config: nil,
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
			got := tt.target.Initialize(tt.config, logger)
			if got != tt.want {
				t.Errorf("Initialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_Finalize(t *testing.T) {
	tests := []struct {
		name   string
		target controller.Controller
	}{
		{
			name:   "Full function",
			target: testTarget["full"],
		},
		{
			name:   "Non init function",
			target: testTarget["non_init"],
		},
		{
			name:   "Non output function",
			target: testTarget["non_output"],
		},
		{
			name:   "Non final function",
			target: testTarget["non_final"],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.target.Finalize()
		})
	}
}

func TestController_Output(t *testing.T) {
	tests := []struct {
		name   string
		target controller.Controller
		output *[]byte
		want   bool
	}{
		{
			name:   "Full function",
			target: testTarget["full"],
			output: nil,
			want:   true,
		},
		{
			name:   "Non init function",
			target: testTarget["non_init"],
			output: nil,
			want:   true,
		},
		{
			name:   "Non output function",
			target: testTarget["non_output"],
			output: nil,
			want:   false,
		},
		{
			name:   "Non final function",
			target: testTarget["non_final"],
			output: nil,
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.Output(tt.output)
			if got != tt.want {
				t.Errorf("Output() = %v, want %v", got, tt.want)
			}
		})
	}
}
