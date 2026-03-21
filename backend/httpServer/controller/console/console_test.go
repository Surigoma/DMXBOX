package console_test

import (
	"backend/httpServer/controller/console"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.bug.st/serial"
)

func TestConsoleAPIv1(t *testing.T) {
	ports, err := serial.GetPortsList()
	if err != nil {
		t.Error("Failed to get ports", "err", err)
	}
	if len(ports) <= 0 {
		t.Skip("Can get not ports")
	}
	gin.SetMode(gin.TestMode)
	engine := gin.Default()
	engine.GET("/v1/console", console.GetConsolesV1)
	t.Run("Can get ports", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/console", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, w.Code, 200)
		body := w.Body.String()
		t.Log("result:", body, "req:", ports)
		var result []string
		err := json.Unmarshal([]byte(body), &result)
		if err != nil {
			t.Error("Failed to unmarshaling.")
		}
		assert.Equal(t, result, ports)
	})
}
