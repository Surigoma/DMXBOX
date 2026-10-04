package dmx_test

import (
	"backend/config"
	dmxserver "backend/dmxServer"
	"backend/httpServer/controller/dmx"
	"backend/packageModule"
	"encoding/json"
	"io"
	"log/slog"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestFadeStateAPI(t *testing.T) {
	config.InitializeConfig()
	c := config.Get()
	c.Dmx.Groups = map[string]config.DMXGroup{"stage": {Devices: []config.DMXDevice{{Model: "dimmer", Channel: 1, MaxValue: []uint{255}}}}}
	module := packageModule.PackageModule{Logger: slog.New(slog.NewTextHandler(io.Discard, nil)), Wg: &sync.WaitGroup{}}
	if !dmxserver.Initialize(&module, &c) {
		t.Fatal("initialize failed")
	}
	defer dmxserver.CleanupDMXServer()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/v1/fade-state", dmx.GetFadeStatesV1)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest("GET", "/v1/fade-state", nil))
	var states map[string]dmxserver.FadeState
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	if err := json.Unmarshal(w.Body.Bytes(), &states); err != nil {
		t.Fatal(err)
	}
	if state, ok := states["stage"]; !ok || state.State != "idle" || state.Level != 0 {
		t.Fatal(states)
	}
}
