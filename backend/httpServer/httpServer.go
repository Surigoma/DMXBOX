package httpServer

import (
	"backend/config"
	"backend/docs"
	"backend/httpServer/controller"
	"backend/message"
	"backend/packageModule"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var wg *sync.WaitGroup
var logger *slog.Logger
var engine *gin.Engine
var listenAddr string

var HttpServer packageModule.PackageModule = packageModule.PackageModule{
	ModuleName:     "http",
	Initialize:     Initialize,
	Run:            StartHTTP,
	MessageHandler: handleMessage,
}

//	@title			DMX BOX
//	@version		0.1
//	@description	For my frendz.
//	@termsOfService	http://swagger.io/terms/

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

func Initialize(module *packageModule.PackageModule, config *config.Config) bool {
	listenAddr = fmt.Sprintf("%s:%d", config.Http.IP, config.Http.Port)
	logger = module.Logger
	engine = registerEndPoints()
	wg = module.Wg
	logger.Info("Hello http server", "addr", "http://"+listenAddr)
	return true
}

func registerEndPoints() *gin.Engine {
	route := gin.Default()
	route.Use(sloggin.New(logger))
	route.Use(gin.Recovery())
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

func handleMessage(mes message.Message) int {
	switch mes.Arg.Action {
	case "stop":
		return -1
	}
	return 0
}

func StartHTTP() {
	wg.Add(1)
	defer wg.Done()
	err := engine.Run(listenAddr)
	if err != nil {
		slog.Error("Failed to setup error", "error", err)
		return
	}
}
