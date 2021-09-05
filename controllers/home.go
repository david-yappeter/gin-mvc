package controllers

import (
	"myapp/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type homeControllerInterface interface {
	IndexPage(*gin.Context)
}

type homeController struct{}

var (
	HomeController homeControllerInterface
)

func init() {
	HomeController = new(homeController)
}

func (controllers *homeController) IndexPage(c *gin.Context) {
	data, err := userData(c)
	if err == service.ErrRecordNotFound {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"title": "Home Page",
		"user":  data,
	})
}
