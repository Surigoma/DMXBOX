package dmxserver

import (
	"backend/config"
	"backend/dmxServer/controller"
	"backend/dmxServer/controller/console"
	ftdi_controller "backend/dmxServer/controller/ftdi"
	device "backend/dmxServer/devices"
	wclight "backend/dmxServer/devices/WCLight"
	"backend/dmxServer/devices/dimmer"
	"backend/dmxServer/fps"
	"backend/message"
	"backend/packageModule"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

var logger *slog.Logger
var wg *sync.WaitGroup
var renderWg sync.WaitGroup
var deviceTypes map[string]func() *device.DMXDevice = make(map[string]func() *device.DMXDevice)
var devices map[uuid.UUID]device.DMXDevice = make(map[uuid.UUID]device.DMXDevice)
var renderTypes map[string]func() *controller.Controller = make(map[string]func() *controller.Controller)
var renderers map[string]*controller.Controller = make(map[string]*controller.Controller)
var rendered []byte = make([]byte, 513)
var fpsController *fps.FPSController
var counter int = 0

type EffectParam struct {
	Duration float32
	Delay    float32
	Fps      float32
}

var param EffectParam = EffectParam{
	Duration: 0.7,
	Delay:    0.0,
}

var DMXServer packageModule.PackageModule = packageModule.PackageModule{
	ModuleName:     "dmx",
	Initialize:     Initialize,
	Run:            StartDMX,
	MessageHandler: handleMessage,
}

func Initialize(module *packageModule.PackageModule, config *config.Config) bool {
	logger = module.Logger
	renderWg = sync.WaitGroup{}
	wg = module.Wg
	param.Duration = config.Dmx.FadeInterval
	param.Delay = config.Dmx.Delay
	param.Fps = config.Dmx.Fps

	deviceTypes["dimmer"] = dimmer.NewDimmer
	deviceTypes["wclight"] = wclight.NewWCLight
	renderTypes["console"] = console.NewConsole
	renderTypes["ftdi"] = ftdi_controller.NewFTDI

	for _, controller := range config.Output {
		AddController(controller, config)
	}
	for _, device := range config.Dmx.Devices {
		AddDevice(device.Model, device.Channel, device.MaxValue)
	}
	return true
}

func handleMessage(mes message.Message) int {
	switch mes.Arg.Action {
	case "stop":
		fpsController.Stop()
		return -1
	}
	return 0
}

func AddDevice(deviceType string, channel uint8, maxValue []uint8) bool {
	generator, ok := deviceTypes[deviceType]
	if !ok {
		logger.Warn("Unsupported type", "type", deviceType)
		return false
	}
	dev := generator()
	if !dev.Initialize(channel, maxValue, &rendered) {
		logger.Error("Failed to initialize device.", "dev", dev)
		return false
	}
	return true
}

func AddController(model string, config *config.Config) bool {
	generator, ok := renderTypes[model]
	if !ok {
		logger.Warn("Unsupported render model", "model", model)
		return false
	}
	if _, ok := renderers[model]; ok {
		return true
	}
	dev := generator()
	if !dev.Initialize(param.Fps, config, logger) {
		logger.Error("Failed to initialize render model.", "model", dev.Model)
		return false
	}
	renderers[model] = dev
	return true
}

func Render() {
	renderWg.Add(len(devices))
	for _, device := range devices {
		device.Update(&renderWg)
	}
	renderWg.Wait()
}

func DMXThread() bool {
	if counter == 0 {
		logger.Debug("fps", "fps", fpsController.GetFPS())
	}
	counter = (counter + 1) % 10
	for _, r := range renderers {
		r.Output(&rendered)
	}
	return true
}

func Finalize() {
	wg.Done()
}

func StartDMX() {
	fpsController = fps.NewFPS(param.Fps, DMXThread, Finalize)
	if fpsController == nil {
		logger.Error("Failed to setup FPS controller.")
		return
	}
	go fpsController.Run()
}
