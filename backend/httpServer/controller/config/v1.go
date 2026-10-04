package config

import (
	"backend/config"
	"backend/message"
	"backend/packageModule"
	"net/http"
	"slices"

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
	Result          bool   `json:"result"`
	Message         string `json:"message"`
	RestartRequired bool   `json:"restartRequired,omitempty"`
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
	manager := packageModule.GetModuleManager()
	var newConfig config.Config
	err := g.ShouldBindJSON(&newConfig)
	if err != nil {
		g.JSON(http.StatusBadRequest, ConfigResult{
			Result:  false,
			Message: err.Error(),
		})
		return
	}
	if err := newConfig.Validate(); err != nil {
		g.JSON(http.StatusBadRequest, ConfigResult{Message: err.Error()})
		return
	}
	oldConfig := config.Get()
	oldInputs, newInputs := slices.Clone(oldConfig.Input.Modules), slices.Clone(newConfig.Input.Modules)
	slices.Sort(oldInputs)
	slices.Sort(newInputs)
	restartRequired := !slices.Equal(oldInputs, newInputs) || slices.Contains(oldConfig.Output.Target, "osc") != slices.Contains(newConfig.Output.Target, "osc")
	if active := manager.GetModules(); len(active) > 0 {
		desired := append(slices.Clone(newInputs), "dmx")
		if slices.Contains(newConfig.Output.Target, "osc") {
			desired = append(desired, "osc")
		}
		slices.Sort(active)
		slices.Sort(desired)
		restartRequired = !slices.Equal(active, desired)
	}
	if err := config.SaveAndSet(newConfig); err != nil {
		g.JSON(http.StatusInternalServerError, ConfigResult{
			Result:  false,
			Message: err.Error(),
		})
		return
	}

	if restartRequired {
		g.JSON(http.StatusOK, ConfigResult{Result: true, RestartRequired: true, Message: "Saved. Restart DMXBOX to apply the module changes and other settings in this update."})
		return
	}

	go func() {
		manager.SendMessageAll(message.Message{
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
