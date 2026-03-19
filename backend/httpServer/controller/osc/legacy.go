package osc

import (
	"backend/message"
	"backend/packageModule"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Mute control
//
//	@Summary	Control a DMX
//	@Schemes
//	@Description	Control a mute status using OSC
//	@Tags			OSC,Legacy
//	@Accept			json
//	@Produce		json
//
//	@Param			mute	query		bool	false	"Mute"
//
//	@Success		200		{object}	OSCResult
//	@Failure		400		{object}	OSCResult
//	@Router			/mute [get]
func LegacyMute(g *gin.Context) {
	arg := map[string]string{}
	isMuteStr := g.Query("mute")
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
	g.JSON(http.StatusOK, OSCResult{
		Result: "OK",
	})
}
