package module

import (
	"backend/config"
	"backend/dmxServer/controller"
	"fmt"
	"log/slog"
)

var counter = 0

func NewConsole() *controller.Controller {
	return &controller.Controller{
		Model: "console",
		ModInitialize: func(*config.Config, *slog.Logger) bool {
			return true
		},
		ModOutput: func(output *[]byte) bool {
			if counter == 0 {
				fmt.Printf("%v\n", *output)
			}
			counter = (counter + 1) % 100
			return true
		},
	}
}
