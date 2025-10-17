package console

import (
	"backend/config"
	"backend/dmxServer/controller"
	"fmt"
	"log/slog"
)

var counter = 0

func NewConsole() *controller.Controller {
	return &controller.Controller{
		Model:         "console",
		ModInitialize: Initialize,
		ModOutput:     Outputfunc,
	}
}

func Initialize(config *config.Config, log *slog.Logger) bool {
	return true
}

func Outputfunc(output *[]byte) bool {
	if counter == 0 {
		fmt.Printf("%v\n", *output)
	}
	counter = (counter + 1) % 100
	return true
}
