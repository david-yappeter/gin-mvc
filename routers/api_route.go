package routers

import (
	"myapp/controllers"

	"github.com/gin-gonic/gin"
)

func restApiRoute(router *gin.Engine) {

	apiRoute := router.Group("/api")
	apiRoute.GET("/", controllers.UserController.CheckToken)

	apiRoute.POST("/users/register", controllers.UserController.Register)
	apiRoute.POST("/users/login", controllers.UserController.Login)
}
