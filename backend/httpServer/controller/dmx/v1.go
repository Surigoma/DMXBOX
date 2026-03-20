package dmx

import (
	dmxserver "backend/dmxServer"
	"backend/message"
	"backend/packageModule"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FadeResult struct {
	Result string         `json:"result"`
	Error  map[string]any `json:"err,omitempty"`
}

// FadeV1 In/Out control
//
//	@Summary	Control a DMX
//	@Schemes
//	@Description	Control a DXM using fade in/out
//	@Tags			DMX,v1
//	@Accept			json
//	@Produce		json
//
//	@Param			group		path		string	true	"Name of DMX group"
//	@Param			isIn		query		bool	false	"is FadeV1 In"
//	@Param			interval	query		int		false	"Onetime Interval"
//	@Param			duration	query		int		false	"Onetime duration"
//
//	@Success		200		{object}	FadeResult
//	@Failure		400		{object}	FadeResult
//	@Failure		500		{object}	FadeResult
//	@Router			/v1/fade/{group} [post]
func FadeV1(g *gin.Context) {
	group := g.Param("group")
	if group == "" {
		g.JSON(http.StatusBadRequest, map[string]any{
			"result": "group is needed.",
		})
		return
	}
	isInStr := g.Query("isIn")
	if isInStr == "" {
		isInStr = "true"
	}
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
		g.JSON(http.StatusInternalServerError, FadeResult{
			Result: "Message send error",
			Error: map[string]any{
				"arg": msg,
			},
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
//	@Tags			DMX,v1
//	@Accept			json
//	@Produce		json
//
//	@Success		200		{object}	map[string]config.DMXGroup
//	@Router			/v1/config/fade [get]
func GetFadeConfigV1(g *gin.Context) {
	config := dmxserver.GetConfig()
	g.JSON(http.StatusOK, config)
}
