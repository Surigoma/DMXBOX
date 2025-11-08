package packageModule

import (
	"backend/config"
	"backend/message"
	"log/slog"
	"sync"
	"time"
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
var running bool

func (mgr *ModuleManagerType) Initialize(log *slog.Logger) bool {
	mgr.logger = log
	mgr.modules = make(map[string]*PackageModule)
	mgr.wg = sync.WaitGroup{}
	running = true
	return true
}

func (mgr *ModuleManagerType) Finalize() {
	running = false
	c := make(chan struct{})
	go func() {
		defer close(c)
		mgr.wg.Wait()
	}()
	select {
	case <-c:
		break
	case <-time.After(3 * time.Second):
		mgr.logger.Error("Failed to wait.", "wg", &mgr.wg)
	}
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

func (mgr *ModuleManagerType) ModuleInitialize(logHandler *slog.Handler) {
	for name, module := range mgr.modules {
		module.Logger = slog.New(*logHandler).With("module", name)
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
		go module.MessageProcess(module.ModuleName, module.MessageHandler)
		go module.Run()
	}
}

func (mgr *ModuleManagerType) SendMessage(msg message.Message) bool {
	module, ok := mgr.modules[msg.To]
	if !ok {
		mgr.logger.Warn("Module not found.", "msg", msg)
		return false
	}
	mgr.logger.Debug("Start send", "to", msg.To, "msg", msg)
	select {
	case module.Channel <- msg:
		mgr.logger.Debug("Send message", "to", msg.To, "msg", msg)
		break
	case <-time.After(time.Duration(1 * time.Second)):
		mgr.logger.Error("message send error", "msg", msg)
		return false
	}
	return true
}

func (mgr *ModuleManagerType) GetModules() []string {
	result := make([]string, 0)
	for k := range mgr.modules {
		result = append(result, k)
	}
	return result
}

func (module *PackageModule) MessageProcess(name string, handler func(msg message.Message) int) {
	module.Logger.Debug("Enter message process.", "module", name)
	for running {
		msg := <-module.Channel
		module.Logger.Debug("Catch message", "module", name, "mes", msg)
		if msg.To == module.ModuleName {
			module.Logger.Debug("Message coming", "module", name, "msg", msg)
			if module.MessageHandler(msg) < 0 {
				break
			}
		} else {
			module.Logger.Error("To is mismatch!", "module", name, "msg", msg)
		}
	}
	module.Logger.Debug("Exit message process.", "module", name)
}
