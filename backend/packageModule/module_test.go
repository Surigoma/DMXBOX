package packageModule

import (
	"backend/config"
	"backend/message"
	"log/slog"
	"testing"
)

func TestBrowserOnlyBlocksTCPControl(t *testing.T) {
	config.InitializeConfig()
	manager := GetModuleManager()
	manager.Initialize(slog.New(slog.NewTextHandler(t.Output(), nil)))
	t.Cleanup(func() {
		manager.SetBrowserOnly(false)
		manager.UnregisterAll()
	})
	target := PackageModule{
		ModuleName:     "dmx",
		Initialize:     func(*PackageModule, *config.Config) bool { return true },
		Run:            func() {},
		Stop:           func() {},
		MessageHandler: func(message.Message) int { return 0 },
	}
	manager.RegisterModule("dmx", &target)
	manager.ModuleInitialize(manager.logger, "test")
	manager.SetBrowserOnly(true)

	control := message.Message{To: "dmx", Arg: message.MessageBody{Action: "fade"}}
	if (&PackageModule{ModuleName: "tcp"}).SendMessage(control) {
		t.Fatal("TCP control was accepted in browser-only mode")
	}
	if !(&PackageModule{ModuleName: "http"}).SendMessage(control) {
		t.Fatal("HTTP control was rejected in browser-only mode")
	}
}
