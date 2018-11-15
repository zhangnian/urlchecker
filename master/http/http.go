package http

import (
	"urlchecker/master/config"

	"fmt"

	"github.com/gin-gonic/gin"
)

func Init() {
	gin.SetMode(gin.ReleaseMode)

	go func() {
		r := gin.Default()

		r.POST("/task/add", handlerTaskAdd)
		r.POST("/task/delete", handlerTaskDelete)
		r.POST("/task/update", handlerTaskUpdate)
		r.GET("/task/list", handlerTaskList)

		r.Run(fmt.Sprintf(":%d", config.G_config.ApiPort))
	}()
}
