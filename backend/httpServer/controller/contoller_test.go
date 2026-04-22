package controller_test

import (
	"backend/config"
	"backend/httpServer"
	"backend/httpServer/controller"
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
		Modules: map[string]bool{
			"http": true,
		},
		Http: config.HttpServer{
			IP:          "127.0.0.1",
			Port:        8080,
			AcceptHosts: []string{"*"},
		},
		Output: config.OutputTargets{
			Target: []string{"console"},
		},
	})
	configData := config.Get()

	module := packageModule.PackageModule{}
	sharedLogger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
	packageModule.ModuleManager.Initialize(sharedLogger)
	h := sharedLogger.Handler()
	packageModule.ModuleManager.ModuleInitialize(&h, "test")
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			module.Logger = slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
			if !httpServer.HttpServer.Initialize(&module, &configData) {
				t.Error("Failed to setup http server")
			}
			gin.SetMode(gin.TestMode)
			t.Parallel()
			engine := httpServer.RegisterEndPoints(&configData.Http, "test")
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
