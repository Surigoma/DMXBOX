package dmxserver

import (
	"backend/config"
	device "backend/dmxServer/devices"
	"backend/message"
	"backend/packageModule"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

var logger *slog.Logger
var wg *sync.WaitGroup
var running bool = false
var deviceTypes map[string]device.DMXDevice = make(map[string]device.DMXDevice)
var devices map[uuid.UUID]device.DMXDevice = make(map[uuid.UUID]device.DMXDevice)
var renderd []byte = make([]byte, 513)

var DMXServer packageModule.PackageModule = packageModule.PackageModule{
	ModuleName:     "dmx",
	Initialize:     Initialize,
	Run:            StartDMX,
	MessageHandler: handleMessage,
}

func Initialize(module *packageModule.PackageModule, config *config.Config) bool {
	logger = module.Logger
	wg = module.Wg
	running = true
	return false
}

func handleMessage(mes message.Message) int {
	switch mes.Arg.Action {
	case "stop":
		running = false
		return -1
	}
	return 0
}

func AddDevice(channel uint8, deviceType string) {
}

func DMXThread() {
	defer wg.Done()
	for running {
		continue
	}
}

func StartDMX() {
	go DMXThread()
}
