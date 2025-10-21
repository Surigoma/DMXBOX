package module

import (
	"backend/artnet"
	"backend/config"
	"backend/dmxServer/controller"
	"log/slog"
)

func NewArtnet() *controller.Controller {
	return &controller.Controller{
		Model:         "artnet",
		ModInitialize: InitializeArtnet,
		ModOutput:     OutputArtnet,
		ModFinalize:   FinalizeArtnet,
	}
}

var a artnet.Artnet

func InitializeArtnet(c *config.Config, log *slog.Logger) bool {
	a = artnet.Artnet{
		TargetAddr: c.Output.Artnet.Address,
	}

	if !a.Initialize(log, c) {
		return false
	}

	return a.Start()
}

func OutputArtnet(data *[]byte) bool {
	a.SendDMXData(data)
	return true
}
func FinalizeArtnet() {
	a.Stop()
}
