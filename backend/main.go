package main

import (
	"backend/config"
	"backend/message"
	"backend/packageModule"
	tcpserver "backend/tcpServer"
	"log/slog"
	"os"
	"os/signal"
	"sync"
)

var channel chan message.Message
var wg sync.WaitGroup
var log *slog.Logger
var logHandler slog.Handler
var modules map[string]*packageModule.PackageModule

func registerModule() {
	modules = make(map[string]*packageModule.PackageModule)
	//modules["http"] = &httpServer.HttpServer
	modules["tcp"] = &tcpserver.TcpServer
}
func moduleInitialize() {
	wg = sync.WaitGroup{}
	channel = make(chan message.Message)
	for name, module := range modules {
		if !module.Initialize(packageModule.PackageModuleParam{
			Logger:  *registerLog(name),
			Config:  config.ConfigData,
			Wg:      &wg,
			Channel: channel,
		}) {
			log.Error("Failed to initialize", "module", name)
		}
	}
}
func moduleStart() {
	for _, module := range modules {
		go module.Run()
	}
}
func handleMessage(mes message.Message) int {
	res := 0
	log.Info("Process message", "mes", mes)
	switch mes.Arg.Action {
	case "stop":
		for k := range modules {
			channel <- message.Message{
				To: k,
				Arg: message.MessageBody{
					Action: "stop",
					Arg:    nil,
				},
			}
		}
		res = -1
	}
	return res
}
func mainProcess() {
	for {
		mes := <-channel
		log.Info("Receive message", "mes", mes)
		if mes.To == "main" {
			if handleMessage(mes) < 0 {
				break
			}
			continue
		}
		channel <- mes
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
	req := slog.Group("base", "module", module)
	return slog.New(logHandler).With(req)
}

func main() {
	logHandler = slog.NewJSONHandler(os.Stdout, nil)
	log = registerLog("main")
	log.Info("Start Main process")
	config.Load(registerLog("config"))
	registerModule()
	moduleInitialize()
	moduleStart()
	go signalProcess()
	mainProcess()
}
