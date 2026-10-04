package dmxserver_test

import (
	"backend/config"
	dmxserver "backend/dmxServer"
	"backend/dmxServer/controller"
	"backend/message"
	"backend/packageModule"
	"io"
	"log/slog"
	"slices"
	"sync"
	"testing"
)

func TestConcurrentFadeRender(t *testing.T) {
	config.InitializeConfig()
	c := config.Get()
	c.Output.Target = []string{"concurrent-test"}
	c.Dmx.Groups = map[string]config.DMXGroup{"stage": {Devices: []config.DMXDevice{{Model: "wclight", Channel: 510, MaxValue: []uint{255, 128, 0}}}}}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	var output []byte
	dmxserver.RenderTypes["concurrent-test"] = func() *controller.Controller {
		return &controller.Controller{
			ModInitialize: func(*config.Config, *slog.Logger) bool { return true },
			ModOutput:     func(data *[]byte) bool { output = slices.Clone(*data); return true },
			ModFinalize:   func() {},
		}
	}
	defer delete(dmxserver.RenderTypes, "concurrent-test")
	var moduleWg sync.WaitGroup
	module := packageModule.PackageModule{Logger: logger, Wg: &moduleWg}
	if !dmxserver.Initialize(&module, &c) {
		t.Fatal("initialize failed")
	}
	moduleWg.Add(1)
	dmxserver.StartDMX()
	var workers sync.WaitGroup
	for worker := range 3 {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for range 100 {
				switch worker {
				case 0:
					dmxserver.DMXServer.MessageHandler(message.Message{Arg: message.MessageBody{Action: "fade", Arg: map[string]string{"id": "stage", "isIn": "true", "duration": "0"}}})
				case 1:
					dmxserver.Render()
				case 2:
					dmxserver.GetConfig()
				}
			}
		}()
	}
	workers.Wait()
	dmxserver.DMXThread()
	dmxserver.StopDMX()
	moduleWg.Wait()
	if len(output) != 512 || !slices.Equal(output[509:], []byte{255, 128, 0}) {
		t.Fatalf("incorrect final DMX output: %v", output)
	}
}
