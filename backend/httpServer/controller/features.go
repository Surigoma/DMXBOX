package controller

import (
	"backend/packageModule"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get features for backend
//
//	@Summary	Get features
//	@Schemes
//	@Description	Get features for backend
//	@Tags			System
//	@Produce		json
//
//	@Success		200		{object}	[]string
//	@Router			/features [get]
func GetFeatures(ctx *gin.Context) {
	manager := packageModule.GetModuleManager()
	ctx.JSON(http.StatusOK, manager.GetModules())
}
