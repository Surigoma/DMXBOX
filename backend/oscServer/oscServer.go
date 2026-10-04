package oscserver

import (
	"backend/config"
	"backend/message"
	"backend/packageModule"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hypebeast/go-osc/osc"
)

var client *osc.Client
var logger *slog.Logger
var wg *sync.WaitGroup
var ip string
var port int
var sendType string

// Zero means no complete mute command has been sent since initialization.
var muteState atomic.Int32

type MuteState struct {
	IsMute *bool `json:"isMute"`
}

func GetMuteState() MuteState {
	state := muteState.Load()
	if state == 0 {
		return MuteState{}
	}
	isMute := state == 1
	return MuteState{IsMute: &isMute}
}

type OSCFormatter struct {
	Base     string
	Type     string
	Inverse  bool
	Channels []uint
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
	values, ok := dataMap[f.Type]
	if !ok {
		return nil, nil
	}
	return result, values[index]
}

var formatter OSCFormatter

var OscServer packageModule.PackageModule = packageModule.PackageModule{
	ModuleName:     "osc",
	Initialize:     Initialize,
	Run:            StartOSC,
	Stop:           func() { wg.Done() },
	MessageHandler: HandleMessage,
}

func Initialize(module *packageModule.PackageModule, config *config.Config) bool {
	muteState.Store(0)
	if config.Output.Osc.Type != "int" && config.Output.Osc.Type != "float" {
		return false
	}
	logger = module.Logger
	wg = module.Wg
	ip = config.Output.Osc.Ip
	port = int(config.Output.Osc.Port)
	sendType = config.Output.Osc.Type
	formatter = OSCFormatter{
		Base:     config.Output.Osc.Format,
		Type:     config.Output.Osc.Type,
		Inverse:  config.Output.Osc.Inverse,
		Channels: config.Output.Osc.Channels,
	}
	client = osc.NewClient(ip, port)
	return true
}

func HandleMessage(mes message.Message) int {
	switch mes.Arg.Action {
	case "reload":
		return 1
	case "stop":
		return -1
	case "mute":
		isMute := true
		if v, ok := mes.Arg.Arg["isMute"]; ok {
			isMute = v == "true"
		}
		addresses, value := formatter.Render(isMute)
		muteState.Store(0)
		sent := len(addresses) > 0
		for _, addr := range addresses {
			p := osc.NewMessage(addr)
			p.Append(value)
			d, _ := p.MarshalBinary()
			logger.Debug("send", "p", p, "b", d)
			err := client.Send(p)
			if err != nil {
				sent = false
				logger.Error("Drop", "err", err)
			}
			time.Sleep(50 * time.Millisecond)
		}
		if sent {
			if isMute {
				muteState.Store(1)
			} else {
				muteState.Store(2)
			}
		}
	}
	return 0
}

func StartOSC() {
}
