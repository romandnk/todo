package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/romandnk/todo/internal/service"
	"github.com/romandnk/todo/pkg/logger"
)

type taskRoutes struct {
	task   service.Task
	logger logger.Logger
}

func newTaskRoutes(g *gin.RouterGroup, task service.Task, logger logger.Logger) {
	r := &taskRoutes{
		task:   task,
		logger: logger,
	}

	g.POST("/", r.CreateTask)
	g.DELETE("/:id", r.DeleteTaskByID)
	g.PATCH("/:id", r.UpdateTaskByID)
	g.GET("/:id", r.GetTaskByID)
}

func (r *taskRoutes) CreateTask(ctx *gin.Context) {

}

func (r *taskRoutes) DeleteTaskByID(ctx *gin.Context) {

}

func (r *taskRoutes) UpdateTaskByID(ctx *gin.Context) {

}

func (r *taskRoutes) GetTaskByID(ctx *gin.Context) {

}
