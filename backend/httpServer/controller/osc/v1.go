package osc

import (
	"backend/message"
	"backend/packageModule"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OSCResult struct {
	Result string         `json:"result"`
	Error  map[string]any `json:"err,omitempty"`
}

// Mute control
//
//	@Summary	Control a OSC
//	@Schemes
//	@Description	Control a mute status using OSC
//	@Tags			OSC,v1
//	@Accept			json
//	@Produce		json
//
//	@Param			isMute	query		bool	false	"Mute"
//
//	@Success		200		{object}	OSCResult
//	@Failure			400		{object}	OSCResult
//	@Router			/v1/mute [post]
func SendOSCV1(g *gin.Context) {
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
	g.JSON(http.StatusOK, OSCResult{
		Result: "OK",
	})
}
