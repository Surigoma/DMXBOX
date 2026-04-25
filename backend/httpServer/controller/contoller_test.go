package controller_test

import (
	"backend/config"
	"backend/httpServer"
	"backend/httpServer/controller"
	"backend/message"
	"backend/packageModule"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAPIResp(t *testing.T) {
	config.InitializeConfig()
	config.Set(config.Config{
		Input: config.InputTargets{
			Modules: []string{"http"},
			Http: config.HttpServer{
				IP:          "127.0.0.1",
				Port:        8080,
				AcceptHosts: []string{"*"},
			},
		},
		Output: config.OutputTargets{
			Target: []string{"console"},
		},
	})
	configData := config.Get()

	module := packageModule.PackageModule{}

	tests := []struct {
		name   string
		path   string
		args   *map[string]string
		method string
		want   func(code int, body string) bool
	}{
		{
			name:   "Version API",
			method: "GET",
			path:   "/api/version",
			want: func(code int, body string) bool {
				var resp controller.VersionInfo = controller.VersionInfo{}
				e := json.Unmarshal([]byte(body), &resp)
				if e != nil {
					return false
				}
				return code == http.StatusOK && resp.Version == "test"
			},
		},
		{
			name:   "Endpoints API",
			method: "GET",
			path:   "/api/endpoints",
			want: func(code int, body string) bool {
				var resp []controller.EndpointInfo = []controller.EndpointInfo{}
				e := json.Unmarshal([]byte(body), &resp)
				if e != nil {
					return false
				}
				return code == http.StatusOK && len(resp) > 0
			},
		},
		{
			name:   "Features API",
			method: "GET",
			path:   "/api/features",
			want: func(code int, body string) bool {
				var resp []string = []string{}
				e := json.Unmarshal([]byte(body), &resp)
				if e != nil {
					return false
				}
				return code == http.StatusOK && len(resp) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sharedLogger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
			packageModule.ModuleManager.Initialize(sharedLogger)
			h := sharedLogger.Handler()
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
			packageModule.ModuleManager.ModuleInitialize(&h, "test")
			packageModule.ModuleManager.ModuleRun()
			defer packageModule.ModuleManager.Finalize()
			module.Logger = sharedLogger
			if !httpServer.HttpServer.Initialize(&module, &configData) {
				t.Error("Failed to setup http server")
			}
			gin.SetMode(gin.TestMode)
			t.Parallel()
			engine := httpServer.RegisterEndPoints(&configData.Input.Http, "test")
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			engine.ServeHTTP(w, req)
			t.Log("request:", tt.method, req.URL.String())
			t.Log("response:", w.Body.String())
			if !tt.want(w.Code, w.Body.String()) {
				t.Error("Failed to check response", "code", w.Code, "body", w.Body.String())
			}
		})
	}
}
