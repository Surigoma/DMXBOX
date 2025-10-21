package controller

import (
	"backend/message"
	"backend/packageModule"
	"net/http"

	"github.com/gin-gonic/gin"
)

//	@BasePath	/api/v1

type OSCResult struct {
	Result string         `json:"result"`
	Error  map[string]any `json:"err",omitempty`
}

// Fade In/Out control
//
//	@Summary	Control a DMX
//	@Schemes
//	@Description	Control a mute status using OSC
//	@Tags			OSC
//	@Accept			json
//	@Produce		json
//
//	@Param			isMute	query		bool	false	"Mute"
//
//	@Success		200		{object}	OSCResult
//	@Success		400		{object}	OSCResult
//	@Router			/mute [post]
func Osc(g *gin.Context) {
	arg := map[string]string{}
	isMuteStr := g.Query("isMute")
	if isMuteStr == "" {
		isMuteStr = "true"
	}
	arg["isMute"] = isMuteStr
	msg := message.Message{
		To: "osc",
		Arg: message.MessageBody{
			Action: "mute",
			Arg: map[string]string{
				"isMute": isMuteStr,
			},
		},
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
