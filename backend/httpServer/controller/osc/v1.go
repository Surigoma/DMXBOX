package osc

import (
	"backend/httpServer/controller"
	"backend/message"
	oscserver "backend/oscServer"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OSCResult struct {
	Result string         `json:"result"`
	Error  map[string]any `json:"err,omitempty"`
}

// GetMuteStateV1 returns the last completely sent OSC mute command.
//
// @Summary Get the last sent OSC mute state
// @Tags OSC,v1
// @Produce json
// @Success 200 {object} oscserver.MuteState
// @Router /v1/mute-state [get]
func GetMuteStateV1(g *gin.Context) {
	g.JSON(http.StatusOK, oscserver.GetMuteState())
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
	ok := controller.SendControl(g, msg)
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
