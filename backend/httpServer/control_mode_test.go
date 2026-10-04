package httpServer_test

import (
	"backend/config"
	"backend/httpServer"
	"backend/packageModule"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestManualControlCORS(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(t.Output(), nil))
	module := &packageModule.PackageModule{Logger: logger}
	if !httpServer.Initialize(module, &config.Config{
		Input: config.InputTargets{Http: config.HttpServer{
			IP: "127.0.0.1", Port: 8080,
			AcceptHosts: []string{"http://192.168.1.10:5173"},
		}},
	}) {
		t.Fatal("HTTP initialization failed")
	}
	engine := httpServer.RegisterEndPoints(&config.HttpServer{
		AcceptHosts: []string{"http://192.168.1.10:5173"},
	}, "test", module)
	for _, path := range []string{"/api/v1/control-mode", "/api/v1/fade/test", "/api/v1/mute"} {
		request := httptest.NewRequest(http.MethodOptions, path, nil)
		request.Header.Set("Origin", "http://192.168.1.10:5173")
		request.Header.Set("Access-Control-Request-Method", "POST")
		request.Header.Set("Access-Control-Request-Headers", "content-type,x-dmxbox-control")
		response := httptest.NewRecorder()
		engine.ServeHTTP(response, request)
		if response.Code != http.StatusNoContent || !strings.Contains(strings.ToLower(response.Header().Get("Access-Control-Allow-Headers")), "x-dmxbox-control") {
			t.Fatalf("manual control preflight rejected for %s: code=%d headers=%v", path, response.Code, response.Header())
		}
	}
}
