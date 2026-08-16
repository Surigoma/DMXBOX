package module

import (
	"backend/config"
	"backend/dmxServer/controller"
	"errors"
	"log/slog"
	"slices"
	"syscall"
	"time"

	"go.bug.st/serial"
)

var loggerFTDI *slog.Logger
var port serial.Port

var target string = ""
var mode = serial.Mode{
	BaudRate: 250000,
	DataBits: 8,
	StopBits: 2,
	Parity:   serial.NoParity,
}

func NewFTDI() *controller.Controller {
	return &controller.Controller{
		ModInitialize: InitializeFTDI,
		ModOutput:     OutputFTDI,
		ModFinalize:   FinalizeFTDI,
		Model:         "FTDI",
	}
}

func InitializeFTDI(config *config.Config, log *slog.Logger) bool {
	loggerFTDI = log
	target = config.Output.FTDI.Port
	ports, err := serial.GetPortsList()
	loggerFTDI.Debug("Found devices", "ports", ports)
	if err != nil {
		loggerFTDI.Error("Failed setup ports.", "err", err)
		return false
	}
	if !slices.Contains(ports, target) {
		loggerFTDI.Error("Port is not found.", "ports", ports)
		return false
	}
	port, err = serial.Open(target, &mode)
	if err != nil {
		loggerFTDI.Error("Failed to open port", "port", target, "err", err)
		return false
	}
	loggerFTDI.Info("Open port.", "port", target)
	return true
}

var zero = []byte{0}

func OutputFTDI(data *[]byte) bool {
	if port == nil {
		port, _ = serial.Open(target, &mode)
		return false
	}
	var err error
	const maxBreakAttempts = 3
	for attempt := 0; attempt < maxBreakAttempts; attempt++ {
		err = port.Break(time.Millisecond)
		if !errors.Is(err, syscall.EINTR) {
			break
		}
		err = nil
	}
	if err != nil {
		loggerFTDI.Error("Failed write data", "err", err)
		port.Close()
		port = nil
		return false
	}
	port.Write(zero)
	port.Write(*data)
	return true
}

func FinalizeFTDI() {
	loggerFTDI.Info("Close port")
	if port != nil {
		port.Close()
	}
}
