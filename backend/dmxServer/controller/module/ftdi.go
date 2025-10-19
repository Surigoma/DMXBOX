package module

import (
	"backend/config"
	"backend/dmxServer/controller"
	"log/slog"
	"slices"
	"time"

	"go.bug.st/serial"
)

var logger *slog.Logger
var port serial.Port

func NewFTDI() *controller.Controller {
	return &controller.Controller{
		ModInitialize: InitializeFTDI,
		ModOutput:     OutputFTDI,
		ModFinalize:   FinalizeFTDI,
		Model:         "FTDI",
	}
}

func InitializeFTDI(config *config.Config, log *slog.Logger) bool {
	logger = log
	target := config.Hardware.Port
	ports, err := serial.GetPortsList()
	logger.Debug("Found devices", "ports", ports)
	if err != nil {
		logger.Error("Failed setup ports.", "err", err)
		return false
	}
	if !slices.Contains(ports, target) {
		logger.Error("Port is not found.", "ports", ports)
		return false
	}
	mode := serial.Mode{
		BaudRate: 250000,
		DataBits: 8,
		StopBits: 2,
		Parity:   serial.NoParity,
	}
	port, err = serial.Open(target, &mode)
	if err != nil {
		logger.Error("Failed to open port", "port", target, "err", err)
		return false
	}
	logger.Info("Open port.", "port", target)
	return true
}

func OutputFTDI(data *[]byte) bool {
	port.Break(time.Duration(1 * time.Millisecond))
	port.Write(*data)
	return false
}

func FinalizeFTDI() {
	logger.Info("Close port")
	port.Close()
}
