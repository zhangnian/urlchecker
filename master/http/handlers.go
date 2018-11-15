package http

import (
	"urlchecker/master"

	"urlchecker/common"

	"github.com/gin-gonic/gin"
)

func handlerTaskAdd(c *gin.Context) {
	var reqBody Task
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    1,
			"message": "参数错误",
			"data":    nil,
		})

		return
	}

	task := common.Task{
		Name:      reqBody.Name,
		Cron:      reqBody.Cron,
		Uri:       reqBody.Uri,
		Method:    reqBody.Method,
		Locations: reqBody.Locations,
	}
	task.NewId()

	err = master.G_taskMgr.SaveTask(&task)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    1,
			"message": "保存任务失败: " + err.Error(),
			"data":    nil,
		})

		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    nil,
	})
}

func handlerTaskDelete(c *gin.Context) {
	var reqBody TaskIds
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    1,
			"message": "参数错误",
			"data":    nil,
		})

		return
	}

	master.G_taskMgr.DeleteTask(reqBody.TaskIds)
}

func handlerTaskUpdate(c *gin.Context) {

}

func handlerTaskList(c *gin.Context) {

}
