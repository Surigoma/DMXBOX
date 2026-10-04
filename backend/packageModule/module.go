package packageModule

import (
	"backend/config"
	"backend/message"
	"backend/operationlog"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type PackageModule struct {
	Initialize     func(module *PackageModule, config *config.Config) bool
	Run            func()
	Stop           func()
	Wg             *sync.WaitGroup
	Channel        chan message.Message
	Logger         *slog.Logger
	ModuleName     string
	MessageHandler func(msg message.Message) int
	Version        string
}

type ModuleManagerType struct {
	modules     map[string]*PackageModule
	logger      *slog.Logger
	wg          sync.WaitGroup
	messageWg   sync.WaitGroup
	lock        sync.Mutex
	browserOnly atomic.Bool
}

var ModuleManager *ModuleManagerType = nil
var running atomic.Bool
var managerOnce = sync.Once{}

func GetModuleManager() *ModuleManagerType {
	managerOnce.Do(func() {
		ModuleManager = &ModuleManagerType{}
	})
	return ModuleManager
}

func (mgr *ModuleManagerType) Initialize(log *slog.Logger) bool {
	mgr.logger = log
	mgr.modules = make(map[string]*PackageModule)
	mgr.wg = sync.WaitGroup{}
	mgr.messageWg = sync.WaitGroup{}
	mgr.lock = sync.Mutex{}
	mgr.browserOnly.Store(false)
	running.Store(true)
	return true
}

func (mgr *ModuleManagerType) Finalize() {
	running.Store(false)
	c := make(chan struct{})
	go func() {
		mgr.wg.Wait()
		mgr.messageWg.Wait()
		defer close(c)
	}()
	select {
	case <-c:
		break
	case <-time.After(3 * time.Second):
		mgr.logger.Error("Failed to wait.", "wg", &mgr.wg)
	}
}

func (mgr *ModuleManagerType) UnregisterAll() {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	for name := range mgr.modules {
		delete(mgr.modules, name)
	}
}

func (mgr *ModuleManagerType) RegisterModule(name string, module *PackageModule) bool {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	_, e := mgr.modules[name]
	if e {
		mgr.logger.Error("Module exists", "name", name, "module", module)
		return false
	}
	mgr.modules[name] = module
	return true
}

func (mgr *ModuleManagerType) ModuleInitialize(log *slog.Logger, version string) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	configData := config.Get()
	for name, module := range mgr.modules {
		module.Logger = log.With("module", name)
		module.Wg = &mgr.wg
		module.Version = version
		module.Channel = make(chan message.Message, 10)
		if !module.Initialize(module, &configData) {
			mgr.logger.Error("Failed to initialize", "module", name)
			delete(mgr.modules, name)
		}
	}
}

func (mgr *ModuleManagerType) ModuleRun() {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	for _, module := range mgr.modules {
		module.Wg.Add(1)
		module.Run()
		mgr.messageWg.Add(1)
		go func() {
			defer mgr.messageWg.Done()
			module.MessageProcess(module.ModuleName, module.MessageHandler)
		}()
	}
}

func (mgr *ModuleManagerType) SendMessageAll(base message.Message) bool {
	for _, m := range mgr.GetModules() {
		msg := base
		msg.To = m
		if !mgr.SendMessage(msg) {
			return false
		}
	}
	return true
}

func (mgr *ModuleManagerType) SendMessage(msg message.Message) bool {
	return mgr.sendMessage(msg, "")
}

func (module *PackageModule) SendMessage(msg message.Message) bool {
	return GetModuleManager().sendMessage(msg, module.ModuleName)
}

func (mgr *ModuleManagerType) sendMessage(msg message.Message, source string) bool {
	if mgr.browserOnly.Load() && source != "http" && (msg.Arg.Action == "fade" || msg.Arg.Action == "mute") {
		mgr.logger.Warn("Control blocked by browser-only mode", "source", source, "action", msg.Arg.Action)
		return false
	}
	mgr.lock.Lock()
	module, ok := mgr.modules[msg.To]
	mgr.lock.Unlock()
	if !ok {
		mgr.logger.Warn("Module not found.", "msg", msg)
		return false
	}
	mgr.logger.Debug("Start send", "to", msg.To, "msg", msg)
	select {
	case module.Channel <- msg:
		if msg.Arg.Action == "fade" || msg.Arg.Action == "mute" {
			if err := operationlog.Record(source, msg.To, msg.Arg.Action, msg.Arg.Arg); err != nil {
				mgr.logger.Error("Failed to write operation log", "err", err)
			}
		}
		mgr.logger.Debug("Send message", "to", msg.To, "msg", msg)
		break
	case <-time.After(time.Duration(1 * time.Second)):
		mgr.logger.Error("message send error", "msg", msg)
		return false
	}
	return true
}

func (mgr *ModuleManagerType) SetBrowserOnly(enabled bool) {
	mgr.browserOnly.Store(enabled)
}

func (mgr *ModuleManagerType) BrowserOnly() bool {
	return mgr.browserOnly.Load()
}

func (mgr *ModuleManagerType) GetModules() []string {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	result := make([]string, 0)
	for k := range mgr.modules {
		result = append(result, k)
	}
	return result
}

func (module *PackageModule) MessageProcess(name string, handler func(msg message.Message) int) {
	module.Logger.Debug("Enter message process.")
	for running.Load() {
		msg := <-module.Channel
		module.Logger.Debug("Catch message", "mes", msg)
		if msg.To == module.ModuleName {
			module.Logger.Debug("Message coming", "msg", msg)
			if res := module.MessageHandler(msg); res < 0 {
				break
			} else if res == 1 {
				module.Logger.Debug("Module reloading", "msg", msg)
				module.Stop()
				configData := config.Get()
				if !module.Initialize(module, &configData) {
					module.Logger.Error("Failed to initialize")
					return
				}
				module.Wg.Add(1)
				module.Run()
			}
		} else {
			module.Logger.Error("To is mismatch!", "msg", msg)
		}
	}
	module.Stop()
	// This code kills test for TCP server.
	// module.Logger.Debug("Exit message process.")
}
