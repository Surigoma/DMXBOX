package oscserver

import (
	"backend/config"
	"backend/message"
	"backend/packageModule"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/hypebeast/go-osc/osc"
)

var client *osc.Client
var logger *slog.Logger
var wg *sync.WaitGroup
var ip string
var port int
var sendType string

type OSCFormatter struct {
	Base     string
	Type     string
	Inverse  bool
	Channels []uint8
}

func (f *OSCFormatter) Render(mute bool) ([]string, any) {
	dataMap := map[string][]any{
		"int":   {int32(0), int32(1)},
		"float": {float32(0), float32(1)},
	}
	result := []string{}
	for _, v := range f.Channels {
		result = append(result, strings.ReplaceAll(f.Base, "{}", fmt.Sprintf("%d", v)))
	}
	index := 0
	if mute != f.Inverse {
		index = 1
	}
	return result, dataMap[f.Type][index]
}

var formatter OSCFormatter

var OscServer packageModule.PackageModule = packageModule.PackageModule{
	ModuleName:     "osc",
	Initialize:     Initialize,
	Run:            StartOSC,
	MessageHandler: handleMessage,
}

func Initialize(module *packageModule.PackageModule, config *config.Config) bool {
	logger = module.Logger
	wg = module.Wg
	ip = config.Osc.Ip
	port = int(config.Osc.Port)
	sendType = config.Osc.Type
	formatter = OSCFormatter{
		Base:     config.Osc.Format,
		Type:     config.Osc.Type,
		Inverse:  config.Osc.Inverse,
		Channels: config.Osc.Channels,
	}
	return true
}

func handleMessage(mes message.Message) int {
	switch mes.Arg.Action {
	case "stop":
		defer wg.Done()
		return -1
	case "mute":
		isMute := true
		if v, ok := mes.Arg.Arg["isMute"]; ok {
			isMute = v == "true"
		}
		addresses, value := formatter.Render(isMute)
		for _, addr := range addresses {
			p := osc.NewMessage(addr)
			p.Append(value)
			d, _ := p.MarshalBinary()
			logger.Debug("send", "p", p, "b", d)
			err := client.Send(p)
			if err != nil {
				logger.Error("Drop", "err", err)
			}
		}
	}
	return 0
}
func StartOSC() {
	client = osc.NewClient(ip, port)
}
