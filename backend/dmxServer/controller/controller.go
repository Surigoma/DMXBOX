package controller

import (
	"backend/config"
	"log/slog"
)

type Controller struct {
	Model         string
	ModInitialize func(*config.Config, *slog.Logger) bool
	ModOutput     func(*[]byte) bool
	ModFinalize   func()
	logger        *slog.Logger
}

func (c *Controller) Initialize(FPS float32, config *config.Config, log *slog.Logger) bool {
	if c.ModInitialize == nil {
		return false
	}
	c.logger = log
	result := c.ModInitialize(config, log.With("module_dmx", c.Model))
	return result
}

func (c *Controller) Finalize() {
	if c.ModFinalize != nil {
		c.ModFinalize()
	}
}

func (c *Controller) Output(output *[]byte) bool {
	if c.ModOutput != nil {
		return c.ModOutput(output)
	} else {
		return false
	}
}
