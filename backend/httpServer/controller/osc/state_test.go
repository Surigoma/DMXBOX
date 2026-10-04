package osc_test

import (
	"backend/httpServer/controller/osc"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMuteStateAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/v1/mute-state", osc.GetMuteStateV1)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/v1/mute-state", nil))
	if response.Code != http.StatusOK || response.Body.String() != `{"isMute":null}` {
		t.Fatalf("unexpected response: %d %s", response.Code, response.Body.String())
	}
}
