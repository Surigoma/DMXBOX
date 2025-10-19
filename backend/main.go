package main

import (
	"backend/config"
	dmxserver "backend/dmxServer"
	"backend/httpServer"
	"backend/message"
	"backend/packageModule"
	tcpserver "backend/tcpServer"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"os/signal"
	"sync"
)

var channel chan message.Message
var wg sync.WaitGroup
var log *slog.Logger
var logHandler slog.Handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
})
var modules map[string]*packageModule.PackageModule = make(map[string]*packageModule.PackageModule)

func registerModule() {
	modules["http"] = &httpServer.HttpServer
	modules["tcp"] = &tcpserver.TcpServer

	for k, v := range config.ConfigData.Modules {
		if !v {
			continue
		}
		if module, ok := modules[k]; ok {
			packageModule.ModuleManager.RegisterModule(k, module)
		} else {
			log.Error(fmt.Sprintf("Module %s is not support.\nSupport modules are %v", k, maps.Keys(modules)))
		}
	}
	packageModule.ModuleManager.RegisterModule("dmx", &dmxserver.DMXServer)
}

func handleMessage(mes message.Message) int {
	res := 0
	log.Debug("Process message", "mes", mes)
	switch mes.Arg.Action {
	case "stop":
		for _, k := range packageModule.ModuleManager.GetModules() {
			log.Debug("Stopping Module", "target", k)
			packageModule.ModuleManager.SendMessage(
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

func registerLog(module string) *slog.Logger {
	return slog.New(logHandler).With("base", module)
}

func main() {
	channel = make(chan message.Message, 10)
	packageModule.ModuleManager.Initialize(registerLog("manager"))
	log = registerLog("main")
	log.Info("Start Main process")
	config.Load(registerLog("config"))
	registerModule()
	packageModule.ModuleManager.ModuleInitialize(&logHandler)
	go signalProcess()
	defer packageModule.ModuleManager.Finalize()
	packageModule.ModuleManager.ModuleRun()
	mainProcess()
}
