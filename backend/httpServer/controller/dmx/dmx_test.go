package dmx_test

import (
	"backend/config"
	dmxserver "backend/dmxServer"
	"backend/httpServer/controller/dmx"
	"backend/message"
	"backend/packageModule"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func Initialize(t *testing.T, msgChan *chan message.Message) *packageModule.PackageModule {
	sharedLogger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
	packageModule.ModuleManager.Initialize(sharedLogger)
	var dummyModule *packageModule.PackageModule
	if msgChan != nil {
		dummyModule = &packageModule.PackageModule{
			MessageHandler: func(msg message.Message) int {
				*msgChan <- msg
				return 0
			},
			Initialize: func(module *packageModule.PackageModule, config *config.Config) bool { return true },
			Run:        func() {},
			Stop:       func() {},
			ModuleName: "dmx",
		}
		packageModule.ModuleManager.RegisterModule("dmx", dummyModule)
	}
	h := sharedLogger.Handler()
	packageModule.ModuleManager.ModuleInitialize(&h)
	packageModule.ModuleManager.ModuleRun()
	return dummyModule
}

func TestAPIDMX(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.Default()
	engine.POST("/v1/fade", dmx.FadeV1)
	engine.POST("/v1/fade/:group", dmx.FadeV1)
	engine.GET("/legacy/in", dmx.FadeInLegacy)
	engine.GET("/legacy/out", dmx.FadeOutLegacy)
	engine.GET("/legacy/addIn", dmx.AddFadeInLegacy)
	engine.GET("/legacy/addOut", dmx.AddFadeOutLegacy)
	engine.GET("/v1/config", dmx.GetFadeConfigV1)

	tests := []struct {
		name   string
		method string
		path   string
		want   message.Message
	}{
		{
			name:   "Fade API v1 (FadeIn) no query",
			method: "POST",
			path:   "/v1/fade/test",
			want: message.Message{
				To: "dmx",
				Arg: message.MessageBody{
					Action: "fade",
					Arg: map[string]string{
						"id":   "test",
						"isIn": "true",
					},
				},
			},
		},
		{
			name:   "Fade API v1 (FadeIn)",
			method: "POST",
			path:   "/v1/fade/test?isIn=true",
			want: message.Message{
				To: "dmx",
				Arg: message.MessageBody{
					Action: "fade",
					Arg: map[string]string{
						"id":   "test",
						"isIn": "true",
					},
				},
			},
		},
		{
			name:   "Fade API v1 (FadeIn) with option",
			method: "POST",
			path:   "/v1/fade/test?isIn=true&interval=0&duration=0",
			want: message.Message{
				To: "dmx",
				Arg: message.MessageBody{
					Action: "fade",
					Arg: map[string]string{
						"id":       "test",
						"isIn":     "true",
						"interval": "0",
						"duration": "0",
					},
				},
			},
		},
		{
			name:   "Fade API v1 (FadeOut)",
			method: "POST",
			path:   "/v1/fade/test?isIn=false",
			want: message.Message{
				To: "dmx",
				Arg: message.MessageBody{
					Action: "fade",
					Arg: map[string]string{
						"id":   "test",
						"isIn": "false",
					},
				},
			},
		},
		{
			name:   "Fade API v1 (FadeOut) with option",
			method: "POST",
			path:   "/v1/fade/test?isIn=false&interval=0&duration=0",
			want: message.Message{
				To: "dmx",
				Arg: message.MessageBody{
					Action: "fade",
					Arg: map[string]string{
						"id":       "test",
						"isIn":     "false",
						"interval": "0",
						"duration": "0",
					},
				},
			},
		},
		{
			name:   "Fade In API Legacy",
			method: "GET",
			path:   "/legacy/in",
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
			name:   "Fade In API Legacy with option",
			method: "GET",
			path:   "/legacy/in?delay=0&interval=1",
			want: message.Message{
				To: "dmx",
				Arg: message.MessageBody{
					Action: "fade",
					Arg: map[string]string{
						"id":       "stg",
						"isIn":     "true",
						"duration": "1",
						"interval": "0",
					},
				},
			},
		},
		{
			name:   "Fade Out API Legacy",
			method: "GET",
			path:   "/legacy/out",
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
			name:   "Fade Out API Legacy with option",
			method: "GET",
			path:   "/legacy/out?delay=0&interval=1",
			want: message.Message{
				To: "dmx",
				Arg: message.MessageBody{
					Action: "fade",
					Arg: map[string]string{
						"id":       "stg",
						"isIn":     "false",
						"duration": "1",
						"interval": "0",
					},
				},
			},
		},
		{
			name:   "Fade Add In API Legacy",
			method: "GET",
			path:   "/legacy/addIn",
			want: message.Message{
				To: "dmx",
				Arg: message.MessageBody{
					Action: "fade",
					Arg: map[string]string{
						"id":   "aud",
						"isIn": "true",
					},
				},
			},
		},
		{
			name:   "Fade Add Out API Legacy",
			method: "GET",
			path:   "/legacy/addOut",
			want: message.Message{
				To: "dmx",
				Arg: message.MessageBody{
					Action: "fade",
					Arg: map[string]string{
						"id":   "aud",
						"isIn": "false",
					},
				},
			},
		},
	}
	noModuleTest := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "v1",
			method: "POST",
			path:   "/v1/fade/test",
		},
		{
			name:   "legacy",
			method: "GET",
			path:   "/legacy/in",
		},
	}
	for _, tt := range noModuleTest {
		t.Run("No DMX module "+tt.name, func(t *testing.T) {
			Initialize(t, nil)
			defer packageModule.ModuleManager.Finalize()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			engine.ServeHTTP(w, req)
			assert.Equal(t, w.Code, 500)
		})
	}

	t.Run("No group", func(t *testing.T) {
		msgChan := make(chan message.Message)
		defer close(msgChan)
		module := Initialize(t, &msgChan)
		defer packageModule.ModuleManager.Finalize()
		defer module.Wg.Done()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/fade", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, w.Code, 400)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgChan := make(chan message.Message)
			defer close(msgChan)
			module := Initialize(t, &msgChan)
			defer packageModule.ModuleManager.Finalize()
			defer module.Wg.Done()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			engine.ServeHTTP(w, req)
			assert.Equal(t, w.Code, 200)
			select {
			case msg := <-msgChan:
				assert.Equal(t, tt.want, msg)
			case <-time.After(1 * time.Second):
				t.Error("Failed to send message")
			}
		})
	}
	t.Run("Get config", func(t *testing.T) {
		group := map[string]config.DMXGroup{
			"test": {
				Name: "test",
				Devices: []config.DMXDevice{
					{
						Model:    "dimmer",
						Channel:  1,
						MaxValue: []uint{255},
					},
				},
			},
		}
		configData := config.Config{
			Output: config.OutputTargets{
				Target: []string{"console"},
			},
			Dmx: config.DMXServer{
				Groups: group,
			},
		}
		dmxserver.DMXServer.Logger = slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))

		dmxserver.Initialize(&dmxserver.DMXServer, &configData)
		defer dmxserver.CleanupDMXServer()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/config", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, w.Code, 200)
		var result map[string]config.DMXGroup
		body := w.Body.String()
		t.Log(body)
		err := json.Unmarshal([]byte(body), &result)
		if err != nil {
			t.Error("Failed to unmarshal", "err", err)
		}
		assert.Equal(t, result, group)
	})
}
