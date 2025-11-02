package config_test

import (
	"backend/config"
	"log/slog"
	"testing"
)

func TestLoadWithPath(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path string
		want bool
	}{
		struct {
			name string
			path string
			want bool
		}{
			name: "Empty",
			path: "./test/data/empty.json",
			want: true,
		},
		struct {
			name string
			path string
			want bool
		}{
			name: "Not Exist",
			path: "./test/data/notExist.json",
			want: false,
		},
		struct {
			name string
			path string
			want bool
		}{
			name: "Error format",
			path: "./test/data/error.json",
			want: false,
		},
	}
	t.Chdir("..")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
			got := config.LoadWithPath(logger, tt.path)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("LoadWithPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []string{
		"Load default config.",
	}
	t.Chdir("..")
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
			if !config.Load(logger) {
				t.Error("Failed to load config.")
			}
		})
	}
}
