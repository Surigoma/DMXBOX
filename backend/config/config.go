package config

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

type DMXHardware struct {
	Port string `json:"port"`
}
type Artnet struct {
	Address     string `json:"addr"`
	Universe    uint8  `json:"universe"`
	SubUniverse uint8  `json:"subuni"`
	Net         uint8  `json:"net"`
}
type HttpServer struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}
type TCPServer struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}
type DMXDevice struct {
	Model    string  `json:"model"`
	Channel  uint8   `json:"channel"`
	MaxValue []uint8 `json:"max"`
}
type DMXServer struct {
	Devices      map[string][]DMXDevice `json:"devices"`
	FadeInterval float32                `json:"fadeInterval"`
	Delay        float32                `json:"delay"`
	Fps          float32                `json:"fps"`
}

type OutputTargets struct {
	Target []string    `json:"target"`
	DMX    DMXHardware `json:"dmx"`
	Artnet Artnet      `json:"artnet"`
}
type OSCServer struct {
	Ip       string  `json:"ip"`
	Port     uint16  `json:"port"`
	Format   string  `json:"format"`
	Type     string  `json:"type"`
	Inverse  bool    `json:"inverse"`
	Channels []uint8 `json:"channels"`
}
type Config struct {
	Modules map[string]bool `json:"modules"`
	Output  OutputTargets   `json:"output"`
	Http    HttpServer      `json:"http"`
	Tcp     TCPServer       `json:"tcp"`
	Dmx     DMXServer       `json:"dmx"`
	Osc     OSCServer       `json:"osc"`
}

var ConfigData Config

func InitializeConfig() {
	ConfigData = Config{
		Modules: map[string]bool{
			"http": false,
			"tcp":  false,
		},
		Output: OutputTargets{
			Target: []string{"console"},
			DMX: DMXHardware{
				Port: "COM1",
			},
			Artnet: Artnet{
				Address:     "2.255.255.255/8",
				Universe:    0,
				SubUniverse: 0,
				Net:         0,
			},
		},
		Http: HttpServer{
			IP:   "127.0.0.1",
			Port: 8000,
		},
		Tcp: TCPServer{
			IP:   "127.0.0.1",
			Port: 50000,
		},
		Dmx: DMXServer{
			Devices:      make(map[string][]DMXDevice),
			FadeInterval: 0.7,
			Delay:        0.0,
			Fps:          0.0,
		},
		Osc: OSCServer{
			Ip:      "127.0.0.1",
			Port:    8765,
			Format:  "/yosc:req/set/MIXER:Current/InCh/Fader/On/{}/1",
			Type:    "int",
			Inverse: true,
			Channels: []uint8{
				1, 2, 3, 4,
			},
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
