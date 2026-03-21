package oscserver_test

import (
	"backend/config"
	"backend/message"
	oscserver "backend/oscServer"
	"backend/packageModule"
	"log/slog"
	"sync"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestOSCModule(t *testing.T) {
	configData := config.Config{
		Osc: config.OSCServer{
			Ip:       "127.0.0.1",
			Port:     3000,
			Format:   "/{}/1",
			Type:     "float",
			Inverse:  false,
			Channels: []uint{1},
		},
	}
	t.Run("Can start", func(t *testing.T) {
		logger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
		packageModule.ModuleManager.Initialize(logger)
		packageModule.ModuleManager.RegisterModule("osc", &oscserver.OscServer)
		handler := logger.Handler()
		packageModule.ModuleManager.ModuleInitialize(&handler)
		packageModule.ModuleManager.ModuleRun()
		packageModule.ModuleManager.SendMessage(message.Message{
			To: "osc",
			Arg: message.MessageBody{
				Action: "stop",
			},
		})
		packageModule.ModuleManager.Finalize()
	})
	tests := []struct {
		name    string
		message message.Message
		want    int
	}{
		{
			name: "Reload",
			message: message.Message{
				To: "osc",
				Arg: message.MessageBody{
					Action: "reload",
				},
			},
			want: 1,
		},
		{
			name: "Stop",
			message: message.Message{
				To: "osc",
				Arg: message.MessageBody{
					Action: "stop",
				},
			},
			want: -1,
		},
		{
			name: "Render",
			message: message.Message{
				To: "osc",
				Arg: message.MessageBody{
					Action: "mute",
				},
			},
			want: 0,
		},
		{
			name: "Render (with Option)",
			message: message.Message{
				To: "osc",
				Arg: message.MessageBody{
					Action: "mute",
					Arg: map[string]string{
						"isMute": "true",
					},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
			module := packageModule.PackageModule{
				Wg:     &sync.WaitGroup{},
				Logger: logger,
			}
			oscserver.Initialize(&module, &configData)
			got := oscserver.HandleMessage(tt.message)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestOSCModuleRender(t *testing.T) {
	type WantType struct {
		path   []string
		result any
	}
	tests := []struct {
		name      string
		isMute    bool
		formatter oscserver.OSCFormatter
		want      WantType
	}{
		{
			name:   "Format float mute",
			isMute: true,
			formatter: oscserver.OSCFormatter{
				Base:     "/{}/1",
				Type:     "float",
				Inverse:  false,
				Channels: []uint{1},
			},
			want: WantType{
				path: []string{
					"/1/1",
				},
				result: float32(1),
			},
		},
		{
			name:   "Format float unmute",
			isMute: false,
			formatter: oscserver.OSCFormatter{
				Base:     "/{}/1",
				Type:     "float",
				Inverse:  false,
				Channels: []uint{1},
			},
			want: WantType{
				path: []string{
					"/1/1",
				},
				result: float32(0),
			},
		},
		{
			name:   "Format float mute inverse",
			isMute: true,
			formatter: oscserver.OSCFormatter{
				Base:     "/{}/1",
				Type:     "float",
				Inverse:  true,
				Channels: []uint{1},
			},
			want: WantType{
				path: []string{
					"/1/1",
				},
				result: float32(0),
			},
		},
		{
			name:   "Format float unmute inverse",
			isMute: false,
			formatter: oscserver.OSCFormatter{
				Base:     "/{}/1",
				Type:     "float",
				Inverse:  true,
				Channels: []uint{1},
			},
			want: WantType{
				path: []string{
					"/1/1",
				},
				result: float32(1),
			},
		},
		{
			name:   "Format int mute",
			isMute: true,
			formatter: oscserver.OSCFormatter{
				Base:     "/{}/1",
				Type:     "int",
				Inverse:  false,
				Channels: []uint{1},
			},
			want: WantType{
				path: []string{
					"/1/1",
				},
				result: int32(1),
			},
		},
		{
			name:   "Format int unmute",
			isMute: false,
			formatter: oscserver.OSCFormatter{
				Base:     "/{}/1",
				Type:     "int",
				Inverse:  false,
				Channels: []uint{1},
			},
			want: WantType{
				path: []string{
					"/1/1",
				},
				result: int32(0),
			},
		},
		{
			name:   "Format int mute inverse",
			isMute: true,
			formatter: oscserver.OSCFormatter{
				Base:     "/{}/1",
				Type:     "int",
				Inverse:  true,
				Channels: []uint{1},
			},
			want: WantType{
				path: []string{
					"/1/1",
				},
				result: int32(0),
			},
		},
		{
			name:   "Format int unmute inverse",
			isMute: false,
			formatter: oscserver.OSCFormatter{
				Base:     "/{}/1",
				Type:     "int",
				Inverse:  true,
				Channels: []uint{1},
			},
			want: WantType{
				path: []string{
					"/1/1",
				},
				result: int32(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			path, value := tt.formatter.Render(tt.isMute)
			assert.Equal(t, tt.want.path, path)
			assert.Equal(t, tt.want.result, value)
		})
	}
}
