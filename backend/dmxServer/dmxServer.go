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
	"math"
	"strconv"
	"sync"
)

var logger *slog.Logger
var wg *sync.WaitGroup
var runningWg *sync.WaitGroup
var renderWg sync.WaitGroup
var DeviceTypes map[string]func() *device.DMXDevice = map[string]func() *device.DMXDevice{
	"dimmer":  deviceSpec.NewDimmer,
	"wclight": deviceSpec.NewWCLight,
}
var RenderTypes map[string]func() *controller.Controller = map[string]func() *controller.Controller{
	"console": ctrlModule.NewConsole,
	"ftdi":    ctrlModule.NewFTDI,
	"artnet":  ctrlModule.NewArtnet,
}
var groups map[string]Group = make(map[string]Group)
var renderers map[string]*controller.Controller = make(map[string]*controller.Controller)
var rendered []byte = make([]byte, 512)
var FpsController *fps.FPSController
var counter int = 0

// Fade, rendering and output all share the same DMX buffer.
var stateMutex sync.Mutex

type Group struct {
	Name    string
	Devices []*device.DMXDevice
}
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
	Stop:           StopDMX,
	MessageHandler: handleMessage,
}

func Initialize(module *packageModule.PackageModule, config *config.Config) bool {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	cleanupDMXServer()
	FpsController = nil
	runningWg = nil
	logger = module.Logger
	renderWg = sync.WaitGroup{}
	wg = module.Wg
	param.Duration = config.Dmx.FadeInterval
	param.Delay = config.Dmx.Delay
	param.Fps = config.Dmx.Fps

	for _, controller := range config.Output.Target {
		if controller == "osc" {
			continue
		}
		if !addController(controller, config) {
			logger.Error("Failed to setup dmx server: unknown controller", "controller", controller)
			return false
		}
	}
	if len(renderers) <= 0 {
		logger.Error("Failed to setup dmx server: controller is none")
		return false
	}
	for name, groupDevices := range config.Dmx.Groups {
		groups[name] = Group{
			Name:    groupDevices.Name,
			Devices: make([]*device.DMXDevice, len(groupDevices.Devices)),
		}
		for i, device := range groupDevices.Devices {
			groups[name].Devices[i] = makeDevice(device.Model, device.Channel, device.MaxValue)
			if groups[name].Devices[i] == nil {
				logger.Error("Failed to setup dmx server: failed to create device", "group", name, "device", device.Model, "index", i)
				return false
			}
		}
	}
	return true
}
func CleanupDMXServer() {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	cleanupDMXServer()
}
func cleanupDMXServer() {
	for k := range groups {
		delete(groups, k)
	}
	for k := range renderers {
		delete(renderers, k)
	}
}

func handleMessage(mes message.Message) int {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	switch mes.Arg.Action {
	case "reload":
		return 1
	case "stop":
		return -1
	case "fade":
		id := mes.Arg.Arg["id"]
		isInStr, ok := mes.Arg.Arg["isIn"]
		isIn := true
		if ok {
			isIn = isInStr == "true"
		}
		duration := float32(-1)
		interval := float32(-1)
		if argStr, ok := mes.Arg.Arg["duration"]; ok {
			conv, err := strconv.ParseFloat(argStr, 32)
			if err == nil && !math.IsNaN(conv) && !math.IsInf(conv, 0) {
				duration = float32(conv)
			}
		}
		if argStr, ok := mes.Arg.Arg["interval"]; ok {
			conv, err := strconv.ParseFloat(argStr, 32)
			if err == nil && !math.IsNaN(conv) && !math.IsInf(conv, 0) {
				interval = float32(conv)
			}
		}
		targetGroup, ok := groups[id]
		if !ok {
			logger.Error("group is not found", "id", id)
			return 0
		}
		logger.Debug("action fade", "fade", isIn)
		for _, d := range targetGroup.Devices {
			d.Fade(isIn, duration, interval)
		}
	}
	return 0
}

func MakeDevice(deviceType string, channel uint16, maxValue []uint) *device.DMXDevice {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	return makeDevice(deviceType, channel, maxValue)
}
func makeDevice(deviceType string, channel uint16, maxValue []uint) *device.DMXDevice {
	generator, ok := DeviceTypes[deviceType]
	if !ok {
		logger.Warn("Unsupported type", "type", deviceType)
		return nil
	}
	dev := generator()
	castMaxValue := make([]uint8, len(maxValue))
	for i, v := range maxValue {
		if v > 255 {
			return nil
		}
		castMaxValue[i] = uint8(v)
	}
	if !dev.Initialize(channel, castMaxValue, &rendered, &param.Duration) {
		logger.Error("Failed to initialize device.", "dev", dev)
		return nil
	}
	logger.Debug("Add device", "dev", dev.Model, "target_ch", dev.Channel)
	return dev
}

func AddController(model string, config *config.Config) bool {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	return addController(model, config)
}
func addController(model string, config *config.Config) bool {
	generator, ok := RenderTypes[model]
	if !ok {
		logger.Warn("Unsupported render model", "model", model)
		return false
	}
	if _, ok := renderers[model]; ok {
		return true
	}
	dev := generator()
	if !dev.Initialize(config, logger) {
		logger.Error("Failed to initialize render model.", "model", dev.Model)
		return false
	}
	renderers[model] = dev
	return true
}

func GetConfig() map[string]config.DMXGroup {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	result := make(map[string]config.DMXGroup, 0)
	for k, v := range groups {
		result[k] = config.DMXGroup{
			Name:    v.Name,
			Devices: make([]config.DMXDevice, len(v.Devices)),
		}
		for i, d := range v.Devices {
			result[k].Devices[i].Channel = d.Channel
			result[k].Devices[i].MaxValue = make([]uint, len(d.MaxValue))
			for ii, m := range d.MaxValue {
				result[k].Devices[i].MaxValue[ii] = uint(m)
			}
			result[k].Devices[i].Model = d.Model
		}
	}
	return result
}

func Render() bool {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	return render()
}
func render() bool {
	result := false
	for _, deviceGroup := range groups {
		for _, device := range deviceGroup.Devices {
			renderWg.Add(1)
			result = device.Update(&renderWg) || result
		}
	}
	renderWg.Wait()
	return result
}

func DMXThread() bool {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	if counter == 0 {
		logger.Debug("fps", "fps", FpsController.GetFPS())
	}
	if render() || counter%10 == 0 {
		for _, r := range renderers {
			r.Output(&rendered)
		}
	}
	counter = (counter + 1) % 50
	return true
}

func Finalize() {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	logger.Debug("Finalize dmx service")
	for k, r := range renderers {
		logger.Debug("Finalize", "k", k)
		r.Finalize()
	}
	cleanupDMXServer()
	wg.Done()
	if runningWg != nil {
		runningWg.Done()
	}
}

func StartDMX() {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	FpsController = fps.NewFPS(param.Fps, DMXThread, Finalize)
	if FpsController == nil {
		logger.Error("Failed to setup FPS controller.")
		wg.Done()
		return
	}
	runningWg = &sync.WaitGroup{}
	runningWg.Add(1)
	go FpsController.Run()
}

func StopDMX() {
	stateMutex.Lock()
	controller, activeWg := FpsController, runningWg
	stateMutex.Unlock()
	if controller != nil {
		controller.Stop()
	}
	if activeWg != nil {
		activeWg.Wait()
	}
}
