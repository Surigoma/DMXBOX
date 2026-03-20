package httpServer_test

import (
	"backend/config"
	"backend/httpServer"
	"backend/httpServer/controller"
	"backend/httpServer/controller/dmx"
	"backend/message"
	"backend/packageModule"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		title  string
		module *packageModule.PackageModule
		config *config.Config
		want   bool
	}{
		{
			title:  "Success",
			module: &packageModule.PackageModule{},
			config: &config.Config{
				Modules: map[string]bool{
					"http": true,
				},
				Http: config.HttpServer{
					IP:          "127.0.0.1",
					Port:        8080,
					AcceptHosts: []string{"*"},
				},
			},
			want: true,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			tt.module.Logger = slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
			o, _ := json.Marshal(tt.config)
			t.Logf("config %v", string(o))
			got := httpServer.HttpServer.Initialize(tt.module, tt.config)
			gin.SetMode(gin.TestMode)
			if !got {
				t.Error("Failed to initialize HTTP Server")
			}
		})
	}
}
func TestViewStaticFile(t *testing.T) {
	module := packageModule.PackageModule{}
	config := config.Config{
		Modules: map[string]bool{
			"http": true,
		},
		Http: config.HttpServer{
			IP:          "127.0.0.1",
			Port:        8080,
			AcceptHosts: []string{"*"},
		},
	}
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "Show Index file",
			path: "/gui/",
			want: want{
				code: 200,
				body: "TEST SUCCESS",
			},
		},
		{
			name: "Show other index file",
			path: "/gui/test.txt",
			want: want{
				code: 200,
				body: "TEST FILE",
			},
		},
		{
			name: "Redirect to /gui/",
			path: "/",
			want: want{
				code: 307,
				body: "",
			},
		},
	}
	t.Chdir("../test/data")
	for _, tt := range tests {
		t.Run("Can Start with static files: "+tt.name, func(t *testing.T) {
			module.Logger = slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
			if !httpServer.HttpServer.Initialize(&module, &config) {
				t.Error("Failed to setup http server")
			}
			gin.SetMode(gin.TestMode)
			engine := httpServer.RegisterEndPoints(&config.Http)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.path, nil)
			req.RequestURI = tt.path
			engine.ServeHTTP(w, req)
			assert.Equal(t, tt.want.code, w.Code)
			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, w.Body.String())
			}
		})
	}
}

func TestAPIResp(t *testing.T) {
	config.InitializeConfig()
	config.Set(config.Config{
		Modules: map[string]bool{
			"http": true,
		},
		Http: config.HttpServer{
			IP:          "127.0.0.1",
			Port:        8080,
			AcceptHosts: []string{"*"},
		},
		Dmx: config.DMXServer{
			Groups: map[string]config.DMXGroup{
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
			},
		},
		Output: config.OutputTargets{
			Target: []string{"console"},
		},
	})
	configData := config.Get()
	configDataB, _ := json.Marshal(configData)

	module := packageModule.PackageModule{}
	sharedLogger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
	packageModule.ModuleManager.Initialize(sharedLogger)
	dummyModule := packageModule.PackageModule{
		MessageHandler: func(msg message.Message) int {
			t.Log(msg)
			return 0
		},
		Initialize: func(module *packageModule.PackageModule, config *config.Config) bool { return true },
		Run:        func() {},
		Stop:       func() {},
		ModuleName: "dummy",
	}
	packageModule.ModuleManager.RegisterModule("dmx", &dummyModule)
	h := sharedLogger.Handler()
	packageModule.ModuleManager.ModuleInitialize(&h)
	packageModule.ModuleManager.ModuleRun()
	defer packageModule.ModuleManager.Finalize()
	tests := []struct {
		name   string
		path   string
		args   *map[string]string
		method string
		want   func(code int, body string) bool
	}{
		{
			name:   "Health API",
			method: "GET",
			path:   "/api/v1/health",
			want: func(code int, body string) bool {
				var resp *controller.HealthResp = &controller.HealthResp{}
				e := json.Unmarshal([]byte(body), resp)
				if e != nil {
					return false
				}
				return resp.Status == "ok"
			},
		},
		{
			name:   "Get config",
			method: "GET",
			path:   "/api/v1/config/all",
			want: func(code int, body string) bool {
				return code == 200 && body == string(configDataB)
			},
		},
		{
			name:   "Can fade",
			method: "POST",
			path:   "/api/v1/fade/test",
			args:   &map[string]string{"isIn": "true"},
			want: func(code int, body string) bool {
				var v = dmx.FadeResult{}
				err := json.Unmarshal([]byte(body), &v)
				if err != nil {
					return false
				}
				return code == 200 && v.Result == "OK"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			module.Logger = slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
			if !httpServer.HttpServer.Initialize(&module, &configData) {
				t.Error("Failed to setup http server")
			}
			gin.SetMode(gin.TestMode)
			t.Parallel()
			engine := httpServer.RegisterEndPoints(&configData.Http)
			w := httptest.NewRecorder()
			args := ""
			if tt.args != nil {
				args += "?"
				joined := []string{}
				for key, value := range *tt.args {
					joined = append(joined, strings.Join([]string{key, value}, "="))
				}
				args += strings.Join(joined, "&")
			}
			req, _ := http.NewRequest(tt.method, tt.path+args, nil)
			engine.ServeHTTP(w, req)
			t.Log("request:", tt.method, req.URL.String())
			t.Log("response:", w.Body.String())
			if !tt.want(w.Code, w.Body.String()) {
				t.Error("Failed to check response", "code", w.Code, "body", w.Body.String())
			}
		})
	}
}
func TestMessageHandle(t *testing.T) {
	tests := []struct {
		name string
		msg  message.Message
		want int
	}{
		{
			name: "Reload",
			msg: message.Message{
				To: "http",
				Arg: message.MessageBody{
					Action: "reload",
				},
			},
			want: 1,
		},
		{
			name: "Stop",
			msg: message.Message{
				To: "http",
				Arg: message.MessageBody{
					Action: "stop",
				},
			},
			want: -1,
		},
		{
			name: "Any",
			msg: message.Message{
				To: "http",
				Arg: message.MessageBody{
					Action: "test",
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run("Message "+tt.name, func(t *testing.T) {
			got := httpServer.HandleMessage(tt.msg)
			t.Log("Result", got, tt.want)
			assert.Equal(t, got, tt.want)
		})
	}
}
