package routers

import (
	"html/template"
	"myapp/controllers"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

func viewRoute(router *gin.Engine) {
	router.HTMLRender = ginview.New(goview.Config{
		Root:         "views",
		Extension:    ".html",
		Master:       "layouts/master",
		Partials:     []string{},
		Funcs:        template.FuncMap{},
		DisableCache: false,
		Delims:       goview.Delims{Left: "{{", Right: "}}"},
	})
	router.Static("/css", "./views/layouts/css")

    router.GET("/", controllers.HomeController.IndexPage)
	router.GET("/login", controllers.UserController.LoginPage)
}
