package health_test

import (
	"backend/httpServer/controller/health"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestHealthAPIv1(t *testing.T) {
	t.Run("Get health API v1", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		engine := gin.Default()
		engine.GET("/health", health.HealthV1)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, w.Code, 200)
		var result health.HealthResp
		err := json.Unmarshal(w.Body.Bytes(), &result)
		if err != nil {
			t.Error("Failed to unmarshaling.")
		}
		assert.Equal(t, result.Status, "ok")
	})
}
