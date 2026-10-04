package config_test

import (
	baseConfig "backend/config"
	"backend/httpServer/controller/config"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
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
	t.Run("Invalid values do not overwrite saved or running settings", func(t *testing.T) {
		baseConfig.InitializeConfig()
		before := baseConfig.Get()
		baseConfig.Save()
		beforeFile, _ := os.ReadFile("config.json")
		invalid := before
		invalid.Dmx.Fps = 0
		payload, _ := json.Marshal(invalid)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/save", bytes.NewReader(payload))
		engine.ServeHTTP(w, req)
		if w.Code != 400 {
			t.Fatal(w.Code, w.Body.String())
		}
		afterFile, _ := os.ReadFile("config.json")
		if !bytes.Equal(beforeFile, afterFile) || !reflect.DeepEqual(before, baseConfig.Get()) {
			t.Fatal("invalid settings replaced valid configuration")
		}
	})
	t.Run("Module changes are saved with an explicit restart requirement", func(t *testing.T) {
		baseConfig.InitializeConfig()
		changed := baseConfig.Get()
		changed.Input.Modules = []string{"http", "tcp"}
		payload, _ := json.Marshal(changed)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, httptest.NewRequest("POST", "/v1/save", bytes.NewReader(payload)))
		var result config.ConfigResult
		json.Unmarshal(w.Body.Bytes(), &result)
		if w.Code != 200 || !result.RestartRequired || result.Message == "" {
			t.Fatal(w.Code, w.Body.String())
		}
	})
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
				base.Input.Modules = []string{"http"}
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
		base.Input.Modules = []string{"http"}
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
