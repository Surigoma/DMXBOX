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
	Error  map[string]any `json:"err",omitempty`
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
//	@Param			group	query		string	true	"Name of DMX group"
//	@Param			isIn	query		bool	false	"is Fade In"
//
//	@Success		200		{object}	FadeResult
//	@Success		400		{object}	FadeResult
//	@Router			/fade [post]
func Fade(g *gin.Context) {
	result := http.StatusOK
	resultArg := FadeResult{
		Result: "OK",
	}
	defer g.JSON(result, resultArg)
	group, ok := g.Params.Get("group")
	if !ok {
		result = http.StatusBadRequest
		resultArg.Error = map[string]any{
			"result": "group is needed.",
		}
		return
	}
	arg := map[string]string{}
	arg["group"] = group
	isInStr, ok := g.Params.Get("isIn")
	if ok {
		arg["isIn"] = isInStr
	}
	msg := message.Message{
		To: "dmx",
		Arg: message.MessageBody{
			Arg: map[string]string{
				"isIn": isInStr,
			},
		},
	}
	ok = packageModule.ModuleManager.SendMessage(msg)
	if !ok {
		result = http.StatusInternalServerError
		resultArg.Error = map[string]any{
			"result": "Message send error",
			"arg":    msg,
		}
	}
}
