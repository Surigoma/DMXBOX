package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type VersionInfo struct {
	Version string
}

// Get version for backend
//
//	@Summary	Get version
//	@Schemes
//	@Description	Get version for backend
//	@Tags			System
//	@Produce		json
//
//	@Success		200		{object}	VersionInfo
//	@Router			/version [get]
func GetVersion(version string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, VersionInfo{
			Version: version,
		})
	}
}
