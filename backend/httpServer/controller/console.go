package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.bug.st/serial"
)

// Get Console Ports
//
//	@Summary	Get all ports for console
//	@Schemes
//	@Description	Get all ports for console
//	@Tags			Config
//	@Accept			json
//	@Produce		json
//
//	@Success		200		{object}	[]string
//	@Failure		500		{object}	map[string]string
//	@Router			/config/console [get]
func GetConsoles(g *gin.Context) {
	ports, err := serial.GetPortsList()
	if err != nil {
		g.JSON(http.StatusInternalServerError, map[string]any{
			"result": "Can not get ",
			"err":    err.Error(),
		})
	}
	g.JSON(http.StatusOK, ports)
}
