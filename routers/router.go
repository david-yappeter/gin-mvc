package routers

import (
	"myapp/service"

	"github.com/gin-gonic/gin"
)

func WebRoute(router *gin.Engine) {
	router.Use(service.AuthMiddleware())
	docRoute(router)
	restApiRoute(router)
	viewRoute(router)
}
