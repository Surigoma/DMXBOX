package httpServer

import (
	"backend/docs"
	"backend/httpServer/controller"
	"backend/packageModule"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HTTPParams struct {
	logger     slog.Logger
	engine     *gin.Engine
	listenAddr string
}

var p HTTPParams
var wg *sync.WaitGroup

var HttpServer packageModule.PackageModule = packageModule.PackageModule{
	Initialize: Initialize,
	Run:        StartHTTP,
}

//	@title			DMX BOX
//	@version		0.1
//	@description	For my frendz.
//	@termsOfService	http://swagger.io/terms/

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

func Initialize(param packageModule.PackageModuleParam) bool {
	p = HTTPParams{
		listenAddr: fmt.Sprintf("0.0.0.0:%d", param.Config.Http.Port),
		logger:     param.Logger,
		engine:     registerEndPoints(),
	}
	wg = param.Wg
	p.logger.Info("Hello http server", "addr", "http://"+p.listenAddr)
	return true
}

func registerEndPoints() *gin.Engine {
	route := gin.Default()
	route.Use(ginLogFormat(&p.logger))
	route.GET("/", controller.HelloWorld)
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := route.Group("/api/v1")
	{
		eg := v1.Group("/")
		{
			eg.GET("/hello", controller.HelloWorld)
		}
	}
	route.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return route
}

func ginLogFormat(logger *slog.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("gin-request",
			slog.String("time", param.TimeStamp.Format(time.RFC3339)),
			slog.Int("status", param.StatusCode),
			slog.String("latency", param.Latency.String()),
			slog.String("client_ip", param.ClientIP),
			slog.String("method", param.Method),
			slog.String("path", param.Path),
			slog.String("error", param.ErrorMessage),
		)
		return ""
	})
}

func StartHTTP() {
	wg.Add(1)
	defer wg.Done()
	err := p.engine.Run(p.listenAddr)
	if err != nil {
		slog.Error("Failed to setup error", "error", err)
		return
	}
}
