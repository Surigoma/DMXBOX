package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EndpointInfo struct {
	Method   string
	Endpoint string
}

// Get all endpoints for control
//
//	@Summary	Get all endpoints
//	@Schemes
//	@Description	Get all endpoints for control
//	@Tags			System
//	@Produce		json
//
//	@Success		200		{object}	[]EndpointInfo
//	@Router			/endpoints [get]
func GetEndpoints(engine *gin.Engine) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		result := []EndpointInfo{}
		for _, route := range engine.Routes() {
			result = append(result, EndpointInfo{
				Method:   route.Method,
				Endpoint: route.Path,
			})
		}
		ctx.JSON(http.StatusOK, result)
	}
}
