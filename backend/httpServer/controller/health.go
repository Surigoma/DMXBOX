package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//	@BasePath	/api/v1

type HealthResp struct {
	Status string    `json:"status" example:"ok"`
	Time   time.Time `json:"time"`
}

// Liveness godoc
//
//	@Summary	liveness probe
//	@Schemes
//	@Description	do ping
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	HealthResp
//	@Router			/health [get]
func Health(g *gin.Context) {
	g.JSON(http.StatusOK, HealthResp{
		Status: "ok",
		Time:   time.Now(),
	})
}
