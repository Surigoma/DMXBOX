package controller

import (
	"backend/config"
	"log/slog"
	"sync"
)

type Controller struct {
	Model         string
	ModInitialize func(*config.Config, *slog.Logger) bool
	ModOutput     func(*[]byte) bool
	ModFinalize   func()
	logger        *slog.Logger
	wg            *sync.WaitGroup
}

func (c *Controller) Initialize(FPS float32, config *config.Config, log *slog.Logger) bool {
	if c.ModInitialize == nil {
		return false
	}
	c.logger = log
	result := c.ModInitialize(config, log)
	return result
}

func (c *Controller) Finalize() {
	if c.ModFinalize != nil {
		c.ModFinalize()
	}
	c.wg.Done()
}

func (c *Controller) Output(output *[]byte) bool {
	if c.ModOutput != nil {
		return c.ModOutput(output)
	} else {
		return false
	}
}
