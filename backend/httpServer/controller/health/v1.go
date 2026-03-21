package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthResp struct {
	Status string    `json:"status" example:"ok"`
	Time   time.Time `json:"time"`
}

// Liveness godoc
//
//	@Summary	liveness probe
//	@Schemes
//	@Description	do ping
//	@Tags			Health,v1
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	HealthResp
//	@Router			/v1/health [get]
func HealthV1(g *gin.Context) {
	g.JSON(http.StatusOK, HealthResp{
		Status: "ok",
		Time:   time.Now(),
	})
}
