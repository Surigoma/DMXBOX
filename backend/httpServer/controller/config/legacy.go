package config

import (
	"backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Save all config
//
//	@Summary	Save all config
//	@Description	Save all config
//	@Tags			Config,Legacy
//	@Accept			json
//	@Produce		json
//
//	@Success		200		{object}	ConfigResult
//	@Failure		500		{object}	ConfigResult
//	@Router			/config/save [post]
func LegacySave(g *gin.Context) {
	ok, err := config.Save()
	if !ok {
		g.JSON(http.StatusInternalServerError, ConfigResult{
			Result:  false,
			Message: err.Error(),
		})
	}
	g.JSON(http.StatusOK, ConfigResult{
		Result: true,
	})
}
