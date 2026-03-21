package config_test

import (
	baseConfig "backend/config"
	"backend/httpServer/controller/config"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestConfigAPIv1(t *testing.T) {
	t.Chdir(t.TempDir())
	gin.SetMode(gin.TestMode)
	engine := gin.Default()
	engine.GET("/v1/get", config.GetConfigV1)
	engine.POST("/v1/save", config.SetConfigV1)
	engine.POST("/legacy/save", config.LegacySave)
	t.Run("Can get current config", func(t *testing.T) {
		baseConfig.InitializeConfig()
		base := baseConfig.Get()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/get", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, w.Code, 200)
		var body baseConfig.Config
		err := json.Unmarshal(w.Body.Bytes(), &body)
		if err != nil {
			t.Error("Failed to unmarshaling.")
		}
		assert.Equal(t, body, base)
	})
	tests := []struct {
		name   string
		config func() string
		want   bool
	}{
		{
			name: "can save",
			config: func() string {
				base := baseConfig.Get()
				base.Modules = map[string]bool{"http": true}
				baseJson, err := json.Marshal(base)
				if err != nil {
					t.Error("Failed to marshaling.")
				}
				return string(baseJson)
			},
			want: true,
		},
		{
			name: "can not save",
			config: func() string {
				return `{"test": "dummy}` // Broken json
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run("Can save current config: "+tt.name, func(t *testing.T) {
			baseConfig.InitializeConfig()
			configData := tt.config()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/v1/save", bytes.NewReader([]byte(configData)))
			engine.ServeHTTP(w, req)
			if tt.want {
				assert.Equal(t, w.Code, 200)
			} else {
				assert.Equal(t, w.Code, 400)
			}
			var body config.ConfigResult
			bodyString := w.Body.String()
			err := json.Unmarshal([]byte(bodyString), &body)
			if err != nil {
				t.Error("Failed to unmarshaling.")
			}
			assert.Equal(t, body.Result, tt.want)
			if tt.want {
				savedDataB, err := os.ReadFile("./config.json")
				if err != nil {
					t.Error("Failed to save.")
				}
				var savedData baseConfig.Config
				err = json.Unmarshal(savedDataB, &savedData)
				if err != nil {
					t.Error("Failed to unmarshaling on saved file.")
				}
				var base baseConfig.Config
				err = json.Unmarshal([]byte(configData), &base)
				if err != nil {
					t.Error("Failed to unmarshaling on test data.")
				}
				assert.Equal(t, savedData, base)
			}
		})
	}
	t.Run("Can save current config", func(t *testing.T) {
		baseConfig.InitializeConfig()
		base := baseConfig.Get()
		base.Modules = map[string]bool{"http": true}
		baseJson, err := json.Marshal(base)
		if err != nil {
			t.Error("Failed to marshaling.")
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/legacy/save", bytes.NewReader(baseJson))
		engine.ServeHTTP(w, req)
		assert.Equal(t, w.Code, 200)
		var body config.ConfigResult
		bodyString := w.Body.String()
		err = json.Unmarshal([]byte(bodyString), &body)
		if err != nil {
			t.Error("Failed to unmarshaling.")
		}
		assert.Equal(t, body.Result, true)
	})
}
