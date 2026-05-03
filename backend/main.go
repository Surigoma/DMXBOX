package main

import (
	"backend/config"
	dmxserver "backend/dmxServer"
	"backend/httpServer"
	"backend/message"
	oscserver "backend/oscServer"
	"backend/packageModule"
	tcpserver "backend/tcpServer"
	_ "embed"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"os/signal"
	"sync"

	"github.com/SladkyCitron/slogcolor"
)

var channel chan message.Message
var wg sync.WaitGroup
var log *slog.Logger
var modules map[string]*packageModule.PackageModule = make(map[string]*packageModule.PackageModule)

//go:embed .version
var Version string

func registerModule() {
	manager := packageModule.GetModuleManager()
	modules["http"] = &httpServer.HttpServer
	modules["tcp"] = &tcpserver.TcpServer

	for _, name := range config.ConfigData.Input.Modules {
		if module, ok := modules[name]; ok {
			manager.RegisterModule(name, module)
		} else {
			log.Error(fmt.Sprintf("Module %s is not support.\nSupport modules are %v", name, maps.Keys(modules)))
		}
	}
	manager.RegisterModule("dmx", &dmxserver.DMXServer)
	for _, target := range config.ConfigData.Output.Target {
		if target == "osc" {
			manager.RegisterModule("osc", &oscserver.OscServer)
		}
	}
}

func handleMessage(mes message.Message) int {
	manager := packageModule.GetModuleManager()
	res := 0
	log.Debug("Process message", "mes", mes)
	switch mes.Arg.Action {
	case "stop":
		for _, k := range manager.GetModules() {
			log.Debug("Stopping Module", "target", k)
			manager.SendMessage(
				message.Message{
					To: k,
					Arg: message.MessageBody{
						Action: "stop",
						Arg:    nil,
					},
				},
			)
		}
		res = -1
	}
	return res
}
func mainProcess() {
	for {
		msg := <-channel
		log.Info("Receive message", "mes", msg)
		if msg.To == "main" {
			if handleMessage(msg) < 0 {
				break
			}
			continue
		}
	}
	log.Info("Stopping main process.")
	wg.Wait()
}
func signalProcess() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	msg := <-quit
	log.Info("Signal interrupt", "msg", msg)
	channel <- message.Message{
		To: "main",
		Arg: message.MessageBody{
			Action: "stop",
			Arg:    nil,
		},
	}
}

func registerLog(module string, handler slog.Handler) *slog.Logger {
	return slog.New(handler).With("base", module)
}

func main() {
	logOption := *slogcolor.DefaultOptions
	logOption.Level = slog.LevelDebug
	logHandler := slogcolor.NewHandler(os.Stdout, &logOption)
	manager := packageModule.GetModuleManager()
	channel = make(chan message.Message, 10)
	manager.Initialize(registerLog("manager", logHandler))
	log = registerLog("main", logHandler)
	log.Info("Start Main process", "version", Version)
	config.Load(registerLog("config", logHandler))
	registerModule()
	manager.ModuleInitialize(slog.New(logHandler), Version)
	go signalProcess()
	defer manager.Finalize()
	manager.ModuleRun()
	mainProcess()
}
