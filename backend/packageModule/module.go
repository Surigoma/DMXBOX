package packageModule

import (
	"backend/config"
	"backend/message"
	"log/slog"
	"sync"
)

type PackageModuleParam struct {
	Logger  slog.Logger
	Config  config.Config
	Wg      *sync.WaitGroup
	Channel chan message.Message
}

type PackageModule struct {
	Initialize func(param PackageModuleParam) bool
	Run        func()
}
