package config

import (
	"backend/config"
	"backend/message"
	"backend/packageModule"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get all config
//
//	@Summary	Get all config
//	@Schemes
//	@Description	Get all config
//	@Tags			Config,v1
//	@Accept			json
//	@Produce		json
//
//	@Success		200		{object}	config.Config
//	@Router			/v1/config/all [get]
func GetConfigV1(g *gin.Context) {
	config := config.Get()
	g.JSON(http.StatusOK, config)
}

type ConfigResult struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

// Set all config
//
//	@Summary	Set all config
//	@Description	Set all config
//	@Tags			Config,v1
//	@Accept			json
//	@Produce		json
//
//	@Param			request		body		config.Config	true	"Configuration data"
//
//	@Success		200		{object}	ConfigResult
//	@Failure		400		{object}	ConfigResult
//	@Failure		500		{object}	ConfigResult
//	@Router			/v1/config/save [post]
func SetConfigV1(g *gin.Context) {
	var newConfig config.Config
	err := g.ShouldBindJSON(&newConfig)
	if err != nil {
		g.JSON(http.StatusBadRequest, ConfigResult{
			Result:  false,
			Message: err.Error(),
		})
		return
	}
	config.Set(newConfig)
	if ok, err := config.Save(); !ok {
		g.JSON(http.StatusInternalServerError, ConfigResult{
			Result:  false,
			Message: err.Error(),
		})
		return
	}

	go func() {
		packageModule.ModuleManager.SendMessageAll(message.Message{
			To: "",
			Arg: message.MessageBody{
				Action: "reload",
				Arg:    nil,
			},
		})
	}()
	g.JSON(http.StatusOK, ConfigResult{
		Result:  true,
		Message: "",
	})
}
