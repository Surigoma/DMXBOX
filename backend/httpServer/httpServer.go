package httpServer

import (
	"backend/config"
	"backend/docs"
	"backend/httpServer/controller"
	"backend/message"
	"backend/packageModule"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var wg *sync.WaitGroup
var logger *slog.Logger
var engine *gin.Engine
var server *http.Server
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

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

func Initialize(module *packageModule.PackageModule, config *config.Config) bool {
	listenAddr = fmt.Sprintf("%s:%d", config.Http.IP, config.Http.Port)
	logger = module.Logger
	engine = registerEndPoints()
	wg = module.Wg
	server = &http.Server{
		Addr:    listenAddr,
		Handler: engine,
	}
	logger.Info("Hello http server", "addr", "http://"+listenAddr)
	return true
}

func registerEndPoints() *gin.Engine {
	route := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.ConfigData.Http.AcceptHosts
	corsConfig.AllowCredentials = true
	route.Use(cors.New(corsConfig))
	route.Use(sloggin.New(logger))
	route.Use(gin.Recovery())
	if info, _ := os.Stat("./static"); info != nil && info.IsDir() {
		route.Any("/", func(g *gin.Context) {
			g.Redirect(307, "/gui/")
		})
		route.NoRoute(func(g *gin.Context) {
			dir, file := path.Split(g.Request.RequestURI)
			ext := filepath.Ext(file)
			if file == "" || ext == "" {
				g.File("./static/index.html")
			} else {
				targetPath := strings.Join(strings.Split(dir, "/")[2:], "/") + file
				g.File("./static/" + targetPath)
			}
		})
	} else {
		logger.Info("Disable file responser.")
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := route.Group("/api/v1")
	{
		eg := v1.Group("/")
		{
			eg.GET("/hello", controller.HelloWorld)
			eg.POST("/fade/:group", controller.Fade)
			eg.POST("/mute", controller.Osc)
			cfg := eg.Group("/config/")
			{
				cfg.GET("/fade", controller.GetFadeConfig)
			}
		}
	}
	route.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return route
}

func handleMessage(mes message.Message) int {
	switch mes.Arg.Action {
	case "stop":
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Failed to stop http server.", "err", err)
		}
		logger.Error("Finalize")
		return -1
	}
	return 0
}

func StartHTTP() {
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Error("Failed to setup error", "error", err)
		return
	}
}
