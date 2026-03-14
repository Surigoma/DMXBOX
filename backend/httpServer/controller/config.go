package controller

import (
	"backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get all config
//
//	@Summary	Get all config
//	@Schemes
//	@Description	Get all config
//	@Tags			Config
//	@Accept			json
//	@Produce		json
//
//	@Success		200		{object}	config.Config
//	@Router			/config/all [get]
func GetConfig(g *gin.Context) {
	config := config.ConfigData
	g.JSON(http.StatusOK, config)
}

// Set all config
//
//	@Summary	Set all config
//	@Schemes
//	@Description	Set all config
//	@Tags			Config
//	@Accept			json
//	@Produce		json
//
//	@Success		200		{object}	config.Config
//	@Router			/config/save [post]
func SetConfig(g *gin.Context) {
	g.String(http.StatusOK, "OK")
}
