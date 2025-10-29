package dmxserver

import (
	"backend/config"
	"backend/dmxServer/controller"
	ctrlModule "backend/dmxServer/controller/module"
	device "backend/dmxServer/devices"
	deviceSpec "backend/dmxServer/devices/spec"
	"backend/dmxServer/fps"
	"backend/message"
	"backend/packageModule"
	"log/slog"
	"sync"
)

var logger *slog.Logger
var wg *sync.WaitGroup
var renderWg sync.WaitGroup
var deviceTypes map[string]func() *device.DMXDevice = make(map[string]func() *device.DMXDevice)
var devices map[string][]*device.DMXDevice = make(map[string][]*device.DMXDevice)
var renderTypes map[string]func() *controller.Controller = make(map[string]func() *controller.Controller)
var renderers map[string]*controller.Controller = make(map[string]*controller.Controller)
var rendered []byte = make([]byte, 512)
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

	deviceTypes["dimmer"] = deviceSpec.NewDimmer
	deviceTypes["wclight"] = deviceSpec.NewWCLight
	renderTypes["console"] = ctrlModule.NewConsole
	renderTypes["ftdi"] = ctrlModule.NewFTDI
	renderTypes["artnet"] = ctrlModule.NewArtnet

	for _, controller := range config.Output.Target {
		AddController(controller, config)
	}
	for name, groupDevices := range config.Dmx.Devices {
		devices[name] = make([]*device.DMXDevice, len(groupDevices))
		for i, device := range groupDevices {
			devices[name][i] = MakeDevice(device.Model, device.Channel, device.MaxValue)
			if devices[name][i] == nil {
				return false
			}
		}
	}
	return true
}

func handleMessage(mes message.Message) int {
	switch mes.Arg.Action {
	case "stop":
		fpsController.Stop()
		return -1
	case "fade":
		id := mes.Arg.Arg["id"]
		isInStr, ok := mes.Arg.Arg["isIn"]
		isIn := true
		if ok {
			isIn = isInStr == "true"
		}
		targetGroup, ok := devices[id]
		if !ok {
			logger.Error("group is not found", "id", id)
			return 0
		}
		logger.Debug("action fade", "fade", isIn)
		for _, d := range targetGroup {
			d.Fade(isIn)
		}
	}
	return 0
}

func MakeDevice(deviceType string, channel uint8, maxValue []uint) *device.DMXDevice {
	generator, ok := deviceTypes[deviceType]
	if !ok {
		logger.Warn("Unsupported type", "type", deviceType)
		return nil
	}
	dev := generator()
	castMaxValue := make([]uint8, len(maxValue))
	for i, v := range maxValue {
		castMaxValue[i] = uint8(v)
	}
	if !dev.Initialize(channel, castMaxValue, &rendered, &param.Duration) {
		logger.Error("Failed to initialize device.", "dev", dev)
		return nil
	}
	logger.Debug("Add device", "dev", dev)
	return dev
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

func GetConfig() map[string][]config.DMXDevice {
	result := make(map[string][]config.DMXDevice, 0)
	for k, v := range devices {
		result[k] = make([]config.DMXDevice, len(v))
		for i, d := range v {
			result[k][i].Channel = d.Channel
			result[k][i].MaxValue = make([]uint, len(d.MaxValue))
			for ii, m := range d.MaxValue {
				result[k][i].MaxValue[ii] = uint(m)
			}
			result[k][i].Model = d.Model
		}
	}
	return result
}

func Render() bool {
	result := false
	for _, deviceGroup := range devices {
		for _, device := range deviceGroup {
			renderWg.Add(1)
			result = device.Update(&renderWg) || result
		}
	}
	renderWg.Wait()
	return result
}

func DMXThread() bool {
	if counter == 0 {
		logger.Debug("fps", "fps", fpsController.GetFPS())
	}
	if Render() || counter%10 == 0 {
		for _, r := range renderers {
			r.Output(&rendered)
		}
	}
	counter = (counter + 1) % 50
	return true
}

func Finalize() {
	logger.Debug("Finalize dmx service")
	for k, r := range renderers {
		logger.Debug("Finalize", "k", k)
		r.Finalize()
	}
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
