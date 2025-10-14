package packageModule

import (
	"backend/config"
	"backend/message"
	"log/slog"
	"sync"
)

type PackageModule struct {
	Initialize     func(module *PackageModule, config *config.Config) bool
	Run            func()
	Wg             *sync.WaitGroup
	Channel        chan message.Message
	Logger         *slog.Logger
	ModuleName     string
	MessageHandler func(msg message.Message) int
}

type ModuleManagerType struct {
	modules map[string]*PackageModule
	logger  *slog.Logger
	wg      sync.WaitGroup
}

var ModuleManager ModuleManagerType = ModuleManagerType{}

func (mgr *ModuleManagerType) Initialize(log *slog.Logger) bool {
	mgr.logger = log
	mgr.modules = make(map[string]*PackageModule)
	mgr.wg = sync.WaitGroup{}
	return true
}

func (mgr *ModuleManagerType) RegisterModule(name string, module *PackageModule) bool {
	_, e := mgr.modules[name]
	if e {
		mgr.logger.Error("Module exists", "name", name, "module", module)
		return false
	}
	mgr.modules[name] = module
	return true
}

func (mgr *ModuleManagerType) ModuleInitialize(logHander *slog.Handler) {
	for name, module := range mgr.modules {
		module.Logger = slog.New(*logHander).With("module", name)
		module.Wg = &mgr.wg
		module.Channel = make(chan message.Message, 10)
		if !module.Initialize(module, &config.ConfigData) {
			mgr.logger.Error("Failed to initialize", "module", name)
		}
	}
}

func (mgr *ModuleManagerType) ModuleRun() {
	for _, module := range mgr.modules {
		module.Wg.Add(1)
		go module.MessageProcess("tcp", module.MessageHandler)
		go module.Run()
	}
}

func (mgr *ModuleManagerType) SendMessage(msg message.Message) bool {
	module, ok := mgr.modules[msg.To]
	if !ok {
		mgr.logger.Warn("Module not found.", "msg", msg)
		return false
	}
	go func() {
		module.Channel <- msg
	}()
	return true
}

func (module *PackageModule) MessageProcess(name string, handler func(msg message.Message) int) {
	module.Logger.Debug("Enter message process.")
	for {
		msg := <-module.Channel
		module.Logger.Debug("Catch message", "mes", msg)
		if msg.To == module.ModuleName {
			module.Logger.Debug("Message coming", "msg", msg)
			if module.MessageHandler(msg) < 0 {
				break
			}
		} else {
			module.Logger.Error("To is mismatch!", "msg", msg)
		}
	}
	module.Logger.Debug("Exit message process.")
}
