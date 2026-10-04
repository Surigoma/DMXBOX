package controller

import (
	"backend/packageModule"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ControlMode struct {
	BrowserOnly bool `json:"browserOnly"`
}

// IsBrowserRequest identifies manual Web UI control, not authenticated clients.
// The explicit marker also works on LAN HTTP, where Fetch Metadata is absent.
func IsBrowserRequest(ctx *gin.Context) bool {
	if ctx.Request == nil {
		return false
	}
	if ctx.GetHeader("X-DMXBOX-Control") == "web-ui" {
		return true
	}
	site := ctx.GetHeader("Sec-Fetch-Site")
	return ctx.GetHeader("Sec-Fetch-Mode") == "cors" && site != "" && site != "none"
}

// GetControlMode returns the active input restriction.
//
//	@Summary	Get control mode
//	@Tags		System,v1
//	@Produce	json
//	@Success	200	{object}	ControlMode
//	@Router		/v1/control-mode [get]
func GetControlMode(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ControlMode{
		BrowserOnly: packageModule.GetModuleManager().BrowserOnly(),
	})
}

// SetControlMode updates the input restriction without reloading modules.
//
//	@Summary	Set control mode
//	@Tags		System,v1
//	@Accept		json
//	@Produce	json
//	@Param		request	body		ControlMode	true	"Control mode"
//	@Success	200		{object}	ControlMode
//	@Failure	400		{object}	map[string]string
//	@Failure	403		{object}	map[string]string
//	@Router		/v1/control-mode [post]
func SetControlMode(ctx *gin.Context) {
	if !IsBrowserRequest(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "control mode can only be changed from the Web UI"})
		return
	}
	var mode ControlMode
	if err := ctx.ShouldBindJSON(&mode); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	packageModule.GetModuleManager().SetBrowserOnly(mode.BrowserOnly)
	ctx.JSON(http.StatusOK, mode)
}

func EnforceBrowserOnly(ctx *gin.Context) {
	if packageModule.GetModuleManager().BrowserOnly() && !IsBrowserRequest(ctx) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "control is restricted to the Web UI"})
		return
	}
	ctx.Next()
}
