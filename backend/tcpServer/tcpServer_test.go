package tcpserver_test

import (
	"backend/config"
	"backend/message"
	"backend/packageModule"
	tcpserver "backend/tcpServer"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestTCPModule(t *testing.T) {
	config.Set(config.Config{
		Modules: map[string]bool{"tcp": true},
		Tcp: config.TCPServer{
			IP:   "127.0.0.1",
			Port: 50000,
		},
		Dmx: config.DMXServer{
			Groups: map[string]config.DMXGroup{
				"stg": {
					Name: "Stage",
					Devices: []config.DMXDevice{
						{
							Model:    "dimmer",
							Channel:  1,
							MaxValue: []uint{255},
						},
					},
				},
				"aud": {
					Name: "Audience",
					Devices: []config.DMXDevice{
						{
							Model:    "dimmer",
							Channel:  2,
							MaxValue: []uint{255},
						},
					},
				},
			},
		},
	})
	t.Run("Can start", func(t *testing.T) {
		logger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
		packageModule.ModuleManager.Initialize(logger)
		packageModule.ModuleManager.RegisterModule("tcp", &tcpserver.TcpServer)
		handler := logger.Handler()
		packageModule.ModuleManager.ModuleInitialize(&handler)
		packageModule.ModuleManager.ModuleRun()
		packageModule.ModuleManager.SendMessage(message.Message{
			To: "tcp",
			Arg: message.MessageBody{
				Action: "stop",
			},
		})
		packageModule.ModuleManager.Finalize()
	})
	t.Run("Can not start when opened other process", func(t *testing.T) {
		listenAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:50000")
		if err != nil {
			t.Error("Failed to open port", "err", err)
		}
		ln, err := net.ListenTCP("tcp", listenAddr)
		if err != nil {
			t.Error("Failed to start a tcp server", "error", err)
			return
		}
		defer ln.Close()
		logger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
		packageModule.ModuleManager.Initialize(logger)
		packageModule.ModuleManager.RegisterModule("tcp", &tcpserver.TcpServer)
		handler := logger.Handler()
		packageModule.ModuleManager.ModuleInitialize(&handler)
		packageModule.ModuleManager.ModuleRun()
		packageModule.ModuleManager.SendMessage(message.Message{
			To: "tcp",
			Arg: message.MessageBody{
				Action: "stop",
			},
		})
		packageModule.ModuleManager.Finalize()
	})
	tests := []struct {
		name string
		msg  message.Message
		want int
	}{
		{
			name: "Reload",
			msg: message.Message{
				To: "tcp",
				Arg: message.MessageBody{
					Action: "reload",
				},
			},
			want: 1,
		},
		{
			name: "Stop",
			msg: message.Message{
				To: "tcp",
				Arg: message.MessageBody{
					Action: "stop",
				},
			},
			want: -1,
		},
		{
			name: "Other",
			msg: message.Message{
				To: "tcp",
				Arg: message.MessageBody{
					Action: "test",
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tcpserver.HandleMessage(tt.msg)
			assert.Equal(t, got, tt.want)
		})
	}
}

func CreateDummyModule(t *testing.T, msgChan *chan message.Message, moduleName string, logger *slog.Logger) *packageModule.PackageModule {
	var dummyModule *packageModule.PackageModule
	dummyModule = &packageModule.PackageModule{
		MessageHandler: func(msg message.Message) int {
			if msg.Arg.Action == "stop" {
				return -1
			}
			logger.Info("Receive message", "msg", msg)
			go func() {
				*msgChan <- msg
			}()
			return 0
		},
		Initialize: func(module *packageModule.PackageModule, config *config.Config) bool { return true },
		Run: func() {
			dummyModule.Logger.Debug("Start dummy")
		},
		Stop: func() {
			dummyModule.Logger.Debug("Stop dummy")
			dummyModule.Wg.Done()
		},
		ModuleName: moduleName,
		Logger:     logger,
	}
	return dummyModule
}
func TestTCPModuleSocket(t *testing.T) {
	config.Set(config.Config{
		Modules: map[string]bool{"tcp": true},
		Tcp: config.TCPServer{
			IP:   "127.0.0.1",
			Port: 50000,
		},
		Dmx: config.DMXServer{
			Groups: map[string]config.DMXGroup{
				"stg": {
					Name: "Stage",
					Devices: []config.DMXDevice{
						{
							Model:    "dimmer",
							Channel:  1,
							MaxValue: []uint{255},
						},
					},
				},
				"aud": {
					Name: "Audience",
					Devices: []config.DMXDevice{
						{
							Model:    "dimmer",
							Channel:  2,
							MaxValue: []uint{255},
						},
					},
				},
			},
		},
	})

	t.Run("Can open tcp port", func(t *testing.T) {
		channel := make(chan message.Message)
		defer close(channel)
		logger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
		packageModule.ModuleManager.Initialize(logger)
		packageModule.ModuleManager.RegisterModule("tcp", &tcpserver.TcpServer)
		handler := logger.Handler()
		packageModule.ModuleManager.ModuleInitialize(&handler)
		defer packageModule.ModuleManager.Finalize()
		packageModule.ModuleManager.ModuleRun()
		defer packageModule.ModuleManager.SendMessageAll(message.Message{
			To: "tcp",
			Arg: message.MessageBody{
				Action: "stop",
			},
		})

		ln, err := net.Dial("tcp", "127.0.0.1:50000")
		if err != nil {
			t.Error("Failed to start a tcp server", "error", err)
			return
		}

		ln.Write([]byte("test"))
		result := make([]byte, 16)
		length, _ := ln.Read(result)
		t.Log(result[:length])
		ln.Close()
		assert.Equal(t, []byte("ack\r\n"), result[:length])
	})
	testDmx := []struct {
		name       string
		action     string
		moduleName string
		want       message.Message
	}{
		{
			name:       "Can send fade In",
			action:     "fadeIn stg",
			moduleName: "dmx",
			want: message.Message{
				To: "dmx",
				Arg: message.MessageBody{
					Action: "fade",
					Arg: map[string]string{
						"id":   "stg",
						"isIn": "true",
					},
				},
			},
		},
		{
			name:       "Can send fade Out",
			action:     "fadeOut stg",
			moduleName: "dmx",
			want: message.Message{
				To: "dmx",
				Arg: message.MessageBody{
					Action: "fade",
					Arg: map[string]string{
						"id":   "stg",
						"isIn": "false",
					},
				},
			},
		},
		{
			name:       "Can send fade In with option",
			action:     "fadeIn stg interval:0,duration:1",
			moduleName: "dmx",
			want: message.Message{
				To: "dmx",
				Arg: message.MessageBody{
					Action: "fade",
					Arg: map[string]string{
						"id":       "stg",
						"isIn":     "true",
						"interval": "0",
						"duration": "1",
					},
				},
			},
		},
		{
			name:       "Can send fade Out with option",
			action:     "fadeOut stg interval:0,duration:1",
			moduleName: "dmx",
			want: message.Message{
				To: "dmx",
				Arg: message.MessageBody{
					Action: "fade",
					Arg: map[string]string{
						"id":       "stg",
						"isIn":     "false",
						"interval": "0",
						"duration": "1",
					},
				},
			},
		},
		{
			name:       "Can send mute",
			action:     "mute",
			moduleName: "osc",
			want: message.Message{
				To: "osc",
				Arg: message.MessageBody{
					Action: "mute",
					Arg: map[string]string{
						"mute": "true",
					},
				},
			},
		},
		{
			name:       "Can send mute with option",
			action:     "mute true",
			moduleName: "osc",
			want: message.Message{
				To: "osc",
				Arg: message.MessageBody{
					Action: "mute",
					Arg: map[string]string{
						"mute": "true",
					},
				},
			},
		},
		{
			name:       "Can send unmute with option",
			action:     "mute false",
			moduleName: "osc",
			want: message.Message{
				To: "osc",
				Arg: message.MessageBody{
					Action: "mute",
					Arg: map[string]string{
						"mute": "false",
					},
				},
			},
		},
	}
	for _, tt := range testDmx {
		t.Run(tt.name, func(t *testing.T) {
			sharedLogger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
			packageModule.ModuleManager.Initialize(sharedLogger)
			defer packageModule.ModuleManager.Finalize()
			msgChan := make(chan message.Message)
			dummyModule := CreateDummyModule(t, &msgChan, tt.moduleName, sharedLogger)
			packageModule.ModuleManager.RegisterModule(tt.moduleName, dummyModule)
			packageModule.ModuleManager.RegisterModule("tcp", &tcpserver.TcpServer)
			h := sharedLogger.Handler()
			packageModule.ModuleManager.ModuleInitialize(&h)
			packageModule.ModuleManager.ModuleRun()
			defer packageModule.ModuleManager.SendMessageAll(message.Message{
				To: "tcp",
				Arg: message.MessageBody{
					Action: "stop",
				},
			})

			ln, err := net.Dial("tcp", "127.0.0.1:50000")
			if err != nil {
				t.Error("Failed to open socket", "error", err)
				return
			}

			ln.Write([]byte(tt.action))
			result := make([]byte, 16)
			length, _ := ln.Read(result)
			ln.Close()
			assert.Equal(t, []byte("ack\r\n"), result[:length])
			select {
			case msg := <-msgChan:
				assert.Equal(t, msg, tt.want)
			case <-time.After(time.Second):
				t.Error("Failed to send message")
			}
		})
	}
}
