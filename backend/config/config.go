package config

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

type Hardware struct {
	ShowDev bool   `json:"show_dev"`
	URL     string `json:"url"`
}
type HttpServer struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}
type TCPServer struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

type Config struct {
	Modules  map[string]bool `json:"modules"`
	Output   []string        `json:"output"`
	Hardware Hardware        `json:"hw"`
	Http     HttpServer      `json:"http"`
	Tcp      TCPServer       `json:"tcp"`
}

var ConfigData Config

func InitializeConfig() {
	ConfigData = Config{
		Modules: map[string]bool{
			"http": false,
			"tcp":  false,
			"dmx":  false,
		},
		Output: []string{"dmx", "artnet"},
		Hardware: Hardware{
			ShowDev: false,
			URL:     "ftdi://ftdi:232:AB0OXCQ4/1",
		},
		Http: HttpServer{
			IP:   "127.0.0.1",
			Port: 8000,
		},
		Tcp: TCPServer{
			IP:   "127.0.0.1",
			Port: 50000,
		},
	}
}

func Load(logger *slog.Logger) bool {
	InitializeConfig()
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		logger.Error("Failed to load a config file", "error", err)
		return false
	}
	defer jsonFile.Close()
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		logger.Error("Failed to read a jcon file", "error", err)
		return false
	}
	json.Unmarshal(jsonData, &ConfigData)
	logger.Info("Decoded", "config", ConfigData)
	return true
}
