package dmx

import (
	"backend/message"
	"backend/packageModule"
	"net/http"

	"github.com/gin-gonic/gin"
)

const BASE_TARGET = "stg"
const ADDITIONAL_TARGET = "aud"

func createMessageTemplate(g *gin.Context, id string, isIn bool) message.Message {
	isInStr := "false"
	if isIn {
		isInStr = "true"
	}
	msg := message.Message{
		To: "dmx",
		Arg: message.MessageBody{
			Action: "fade",
			Arg: map[string]string{
				"id":   id,
				"isIn": isInStr,
			},
		},
	}
	if intStr := g.Query("interval"); intStr != "" {
		msg.Arg.Arg["duration"] = intStr
	}
	if intStr := g.Query("delay"); intStr != "" {
		msg.Arg.Arg["interval"] = intStr
	}
	return msg
}
func sendMessage(g *gin.Context, msg message.Message) {
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

// Fade In control for old APIs
//
//	@Summary	Control a DMX
//	@Schemes
//	@Description	Control a DXM using fade in for old APIs
//	@Tags			DMX,Legacy
//	@Accept			json
//	@Produce		json
//
//	@Param			delay		query		int		false	"Onetime Interval"
//	@Param			interval	query		int		false	"Onetime duration"
//
//	@Success		200		{object}	FadeResult
//	@Failure		400		{object}	FadeResult
//	@Failure		500		{object}	FadeResult
//	@Router			/fadeIn [get]
func FadeInLegacy(g *gin.Context) {
	msg := createMessageTemplate(g, BASE_TARGET, true)
	sendMessage(g, msg)
}

// Fade Out control for old APIs
//
//	@Summary	Control a DMX
//	@Schemes
//	@Description	Control a DXM using fade out for old APIs
//	@Tags			DMX,Legacy
//	@Accept			json
//	@Produce		json
//
//	@Param			delay		query		int		false	"Onetime Interval"
//	@Param			interval	query		int		false	"Onetime duration"
//
//	@Success		200		{object}	FadeResult
//	@Failure		400		{object}	FadeResult
//	@Failure		500		{object}	FadeResult
//	@Router			/fadeOut [get]
func FadeOutLegacy(g *gin.Context) {
	msg := createMessageTemplate(g, BASE_TARGET, false)
	sendMessage(g, msg)
}

// Fade In control for old APIs
//
//	@Summary	Control a DMX
//	@Schemes
//	@Description	Control a DXM using fade in for old APIs
//	@Tags			DMX,Legacy
//	@Accept			json
//	@Produce		json
//
//	@Param			delay		query		int		false	"Onetime Interval"
//	@Param			interval	query		int		false	"Onetime duration"
//
//	@Success		200		{object}	FadeResult
//	@Failure		400		{object}	FadeResult
//	@Failure		500		{object}	FadeResult
//	@Router			/fadeAddIn [get]
func AddFadeInLegacy(g *gin.Context) {
	msg := createMessageTemplate(g, ADDITIONAL_TARGET, true)
	sendMessage(g, msg)
}

// Additional Fade Out control for old APIs
//
//	@Summary	Control a DMX
//	@Schemes
//	@Description	Control a DXM using fade out for old APIs
//	@Tags			DMX,Legacy
//	@Accept			json
//	@Produce		json
//
//	@Param			delay		query		int		false	"Onetime Interval"
//	@Param			interval	query		int		false	"Onetime duration"
//
//	@Success		200		{object}	FadeResult
//	@Failure		400		{object}	FadeResult
//	@Failure		500		{object}	FadeResult
//	@Router			/fadeAddOut [get]
func AddFadeOutLegacy(g *gin.Context) {
	msg := createMessageTemplate(g, ADDITIONAL_TARGET, false)
	sendMessage(g, msg)
}
