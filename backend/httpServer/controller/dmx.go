package controller

import (
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
//	@Param			group	path		string	true	"Name of DMX group"
//	@Param			isIn	query		bool	false	"is Fade In"
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
