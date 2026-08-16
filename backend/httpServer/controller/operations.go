package controller

import (
	"backend/operationlog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetOperations returns recent control operations.
//
//	@Summary	Get operation log
//	@Tags		System
//	@Produce	json
//	@Param		limit	query		int	false	"Maximum number of entries"
//	@Success	200		{object}	[]operationlog.Entry
//	@Router		/v1/operations [get]
func GetOperations(ctx *gin.Context) {
	limit := 100
	if value, err := strconv.Atoi(ctx.DefaultQuery("limit", "100")); err == nil && value > 0 && value <= 1000 {
		limit = value
	}
	ctx.JSON(http.StatusOK, operationlog.List(limit))
}
