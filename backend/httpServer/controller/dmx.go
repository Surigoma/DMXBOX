package controller

import (
	dmxserver "backend/dmxServer"
	"backend/message"
	"backend/packageModule"
	"net/http"

	"github.com/gin-gonic/gin"
)

//	@BasePath	/api/v1

type FadeResult struct {
	Result string         `json:"result"`
	Error  map[string]any `json:"err,omitempty"`
}

// Fade In/Out control
//
//	@Summary	Control a DMX
//	@Schemes
//	@Description	Control a DXM using fade in/out
//	@Tags			DMX
//	@Accept			json
//	@Produce		json
//
//	@Param			group		path		string	true	"Name of DMX group"
//	@Param			isIn		query		bool	false	"is Fade In"
//	@Param			interval	query		int		false	"Onetime Interval"
//	@Param			duration	query		int		false	"Onetime duration"
//
//	@Success		200		{object}	FadeResult
//	@Success		400		{object}	FadeResult
//	@Router			/fade/{group} [post]
func Fade(g *gin.Context) {
	group := g.Param("group")
	if group == "" {
		g.JSON(http.StatusBadRequest, map[string]any{
			"result": "group is needed.",
		})
		return
	}
	arg := map[string]string{}
	arg["group"] = group
	isInStr := g.Query("isIn")
	if isInStr == "" {
		isInStr = "true"
	}
	arg["isIn"] = isInStr
	msg := message.Message{
		To: "dmx",
		Arg: message.MessageBody{
			Action: "fade",
			Arg: map[string]string{
				"id":   group,
				"isIn": isInStr,
			},
		},
	}
	if intStr := g.Query("interval"); intStr != "" {
		msg.Arg.Arg["interval"] = intStr
	}
	if intStr := g.Query("duration"); intStr != "" {
		msg.Arg.Arg["duration"] = intStr
	}
	ok := packageModule.ModuleManager.SendMessage(msg)
	if !ok {
		g.JSON(http.StatusInternalServerError, map[string]any{
			"result": "Message send error",
			"arg":    msg,
		})
	}
	g.JSON(http.StatusOK, FadeResult{
		Result: "OK",
	})
}

// Get fade config
//
//	@Summary	Control a DMX
//	@Schemes
//	@Description	Control a DXM using fade in/out
//	@Tags			DMX
//	@Accept			json
//	@Produce		json
//
//	@Success		200		{object}	map[string][]config.DMXDevice
//	@Router			/config/fade [get]
func GetFadeConfig(g *gin.Context) {
	config := dmxserver.GetConfig()
	g.JSON(http.StatusOK, config)
}
