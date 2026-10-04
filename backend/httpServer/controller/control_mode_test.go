package controller

import (
	"backend/packageModule"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBrowserOnlyMode(t *testing.T) {
	manager := packageModule.GetModuleManager()
	manager.Initialize(slog.New(slog.NewTextHandler(t.Output(), nil)))
	t.Cleanup(func() { manager.SetBrowserOnly(false) })

	restResponse := httptest.NewRecorder()
	restContext, _ := gin.CreateTestContext(restResponse)
	restContext.Request = httptest.NewRequest("POST", "/api/v1/control-mode", strings.NewReader(`{"browserOnly":true}`))
	restContext.Request.Header.Set("Content-Type", "application/json")
	SetControlMode(restContext)
	if restResponse.Code != http.StatusForbidden || manager.BrowserOnly() {
		t.Fatalf("REST client changed mode: code=%d body=%s", restResponse.Code, restResponse.Body.String())
	}

	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = httptest.NewRequest("POST", "/api/v1/control-mode", strings.NewReader(`{"browserOnly":true}`))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Header.Set("Sec-Fetch-Mode", "cors")
	ctx.Request.Header.Set("Sec-Fetch-Site", "same-origin")
	SetControlMode(ctx)

	if response.Code != http.StatusOK || !manager.BrowserOnly() {
		t.Fatalf("mode was not enabled: code=%d body=%s", response.Code, response.Body.String())
	}

	engine := gin.New()
	engine.GET("/control", EnforceBrowserOnly, func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})
	restResponse = httptest.NewRecorder()
	engine.ServeHTTP(restResponse, httptest.NewRequest(http.MethodGet, "/control", nil))
	if restResponse.Code != http.StatusForbidden {
		t.Fatalf("REST control was accepted: code=%d", restResponse.Code)
	}
	browserResponse := httptest.NewRecorder()
	browserRequest := httptest.NewRequest(http.MethodGet, "/control", nil)
	browserRequest.Header.Set("Sec-Fetch-Mode", "cors")
	browserRequest.Header.Set("Sec-Fetch-Site", "same-origin")
	engine.ServeHTTP(browserResponse, browserRequest)
	if browserResponse.Code != http.StatusNoContent {
		t.Fatalf("Web UI control was rejected: code=%d", browserResponse.Code)
	}
	// LAN HTTP requests have no Sec-Fetch-* headers.
	manualRequest := httptest.NewRequest(http.MethodGet, "/control", nil)
	manualRequest.Header.Set("X-DMXBOX-Control", "web-ui")
	manualResponse := httptest.NewRecorder()
	engine.ServeHTTP(manualResponse, manualRequest)
	if manualResponse.Code != http.StatusNoContent {
		t.Fatalf("LAN Web UI control was rejected: code=%d", manualResponse.Code)
	}
	modeResponse := httptest.NewRecorder()
	modeContext, _ := gin.CreateTestContext(modeResponse)
	modeContext.Request = httptest.NewRequest(http.MethodPost, "/api/v1/control-mode", strings.NewReader(`{"browserOnly":false}`))
	modeContext.Request.Header.Set("Content-Type", "application/json")
	modeContext.Request.Header.Set("X-DMXBOX-Control", "web-ui")
	SetControlMode(modeContext)
	if modeResponse.Code != http.StatusOK || manager.BrowserOnly() {
		t.Fatalf("LAN Web UI could not disable mode: code=%d", modeResponse.Code)
	}
	unrestrictedResponse := httptest.NewRecorder()
	engine.ServeHTTP(unrestrictedResponse, httptest.NewRequest(http.MethodGet, "/control", nil))
	if unrestrictedResponse.Code != http.StatusNoContent {
		t.Fatalf("REST control remained blocked after disabling mode: code=%d", unrestrictedResponse.Code)
	}
}
