package controller

import (
	"backend/message"
	"backend/packageModule"

	"github.com/gin-gonic/gin"
)

const controlModuleKey = "controlModule"

func SetControlModule(module *packageModule.PackageModule) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(controlModuleKey, module)
		ctx.Next()
	}
}

func SendControl(ctx *gin.Context, msg message.Message) bool {
	if module, ok := ctx.Get(controlModuleKey); ok {
		return module.(*packageModule.PackageModule).SendMessage(msg)
	}
	return packageModule.GetModuleManager().SendMessage(msg)
}
