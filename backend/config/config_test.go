package config_test

import (
	"backend/config"
	"encoding/json"
	"log/slog"
	"os"
	"reflect"
	"sync"
	"testing"
)

func TestLoadWithPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "Empty",
			path: "./test/data/empty.json",
			want: true,
		},
		{
			name: "Not Exist",
			path: "./test/data/notExist.json",
			want: false,
		},
		{
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
			if got != tt.want {
				t.Errorf("LoadWithPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	t.Chdir("..")
	t.Run("Load default config.", func(t *testing.T) {
		logger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
		if !config.Load(logger) {
			t.Error("Failed to load config.")
		}
	})
}

func TestGetSet(t *testing.T) {
	t.Run("Equaled Get", func(t *testing.T) {
		config.InitializeConfig()
		got := config.Get()
		if !reflect.DeepEqual(config.ConfigData, got) {
			t.Error("Miss match")
		}
	})
	t.Run("Set", func(t *testing.T) {
		config.InitializeConfig()
		got := config.Get()
		if !reflect.DeepEqual(config.ConfigData, got) {
			t.Error("Miss match")
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			got.Modules["http"] = !got.Modules["http"]
			config.Set(got)
			wg.Done()
		}()
		wg.Wait()
		if !reflect.DeepEqual(config.ConfigData, got) {
			t.Error("Miss match")
		}
	})
}

func TestSave(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Run("Can save config", func(t *testing.T) {
		config.InitializeConfig()
		config.Save()
		wroteData, err := os.ReadFile("./config.json")
		if err != nil {
			t.Error("Failed to load config.json", "err", err)
		}
		var wroteJSON config.Config
		err = json.Unmarshal(wroteData, &wroteJSON)
		if err != nil {
			t.Error("Wrote data is broken.", "err", err)
		}
		if !reflect.DeepEqual(config.ConfigData, wroteJSON) {
			t.Error("Miss match", "orig", config.ConfigData, "wrote", wroteData)
		}
	})
}
