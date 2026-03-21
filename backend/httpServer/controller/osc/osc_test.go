package osc_test

import (
	"backend/config"
	"backend/httpServer/controller/osc"
	"backend/message"
	"backend/packageModule"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func Initialize(t *testing.T, msgChan *chan message.Message) *packageModule.PackageModule {
	sharedLogger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{Level: slog.LevelDebug}))
	packageModule.ModuleManager.Initialize(sharedLogger)
	var dummyModule *packageModule.PackageModule
	if msgChan != nil {
		dummyModule = &packageModule.PackageModule{
			MessageHandler: func(msg message.Message) int {
				*msgChan <- msg
				return 0
			},
			Initialize: func(module *packageModule.PackageModule, config *config.Config) bool { return true },
			Run:        func() {},
			Stop:       func() {},
			ModuleName: "osc",
		}
		packageModule.ModuleManager.RegisterModule("osc", dummyModule)
	}
	h := sharedLogger.Handler()
	packageModule.ModuleManager.ModuleInitialize(&h)
	packageModule.ModuleManager.ModuleRun()
	return dummyModule
}

func TestOSCAPIv1(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		want   message.Message
	}{}
	for _, v := range []string{"/v1/osc", "/legacy/osc"} {
		legStr := ""
		method := "POST"
		args := "?isMute="
		if strings.Contains(v, "/legacy/") {
			legStr = " for Legacy API"
			method = "GET"
			args = "?mute="
		}

		tests = append(tests, struct {
			name   string
			method string
			path   string
			want   message.Message
		}{
			name:   "Can send message" + legStr + " (no option)",
			method: method,
			path:   v,
			want: message.Message{
				To: "osc",
				Arg: message.MessageBody{
					Action: "mute",
					Arg: map[string]string{
						"isMute": "true",
					},
				},
			},
		})
		for _, vv := range []string{"true", "false"} {
			tests = append(tests, struct {
				name   string
				method string
				path   string
				want   message.Message
			}{
				name:   "Can send message" + legStr + " isMute=" + vv,
				method: method,
				path:   v + args + vv,
				want: message.Message{
					To: "osc",
					Arg: message.MessageBody{
						Action: "mute",
						Arg: map[string]string{
							"isMute": vv,
						},
					},
				},
			})
		}
	}
	gin.SetMode(gin.TestMode)
	engine := gin.Default()
	engine.POST("/v1/osc", osc.SendOSCV1)
	engine.GET("/legacy/osc", osc.LegacyMute)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgChan := make(chan message.Message)
			defer close(msgChan)
			Initialize(t, &msgChan)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			engine.ServeHTTP(w, req)
			assert.Equal(t, w.Code, 200)
			select {
			case msg := <-msgChan:
				assert.Equal(t, msg, tt.want)
			case <-time.After(1 * time.Second):
				t.Error("Failed to get message")
			}
		})
	}
	testsCantSend := []struct {
		name   string
		method string
		path   string
		want   message.Message
	}{}
	for _, v := range []string{"/v1/osc", "/legacy/osc"} {
		legStr := ""
		method := "POST"
		if strings.Contains(v, "/legacy/") {
			legStr = " for Legacy API"
			method = "GET"
		}
		testsCantSend = append(testsCantSend, struct {
			name   string
			method string
			path   string
			want   message.Message
		}{
			name:   "Can not send message" + legStr + " (no option)",
			method: method,
			path:   v,
			want: message.Message{
				To: "osc",
				Arg: message.MessageBody{
					Action: "mute",
					Arg: map[string]string{
						"isMute": "true",
					},
				},
			},
		})
	}
	for _, tt := range testsCantSend {
		t.Run(tt.name, func(t *testing.T) {
			Initialize(t, nil)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			engine.ServeHTTP(w, req)
			assert.Equal(t, w.Code, 500)
		})
	}
}
