package ftdi_controller

import (
	"backend/config"
	"backend/dmxServer/controller"
	"log/slog"

	"periph.io/x/host/v3"
	"periph.io/x/host/v3/ftdi"
)

var logger *slog.Logger
var device ftdi.Dev

func NewFTDI() *controller.Controller {
	return &controller.Controller{
		ModInitialize: Initialize,
		ModOutput:     OutputFunc,
		Model:         "FTDI",
	}
}

func Initialize(config *config.Config, log *slog.Logger) bool {
	logger = log
	if _, err := host.Init(); err != nil {
		logger.Error("Failed to setup FTDI controller", "err", err)
		return false
	}
	all := ftdi.All()
	if len(all) == 0 {
		logger.Error("FTDI device is not found.", "devs", all)
		return false
	}
	for _, v := range all {
		log.Debug(v.String())
	}
	return false
}

func OutputFunc(data *[]byte) bool {
	return false
}
