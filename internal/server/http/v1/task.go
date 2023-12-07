package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/romandnk/todo/internal/constant"
	"github.com/romandnk/todo/internal/service"
	"github.com/romandnk/todo/pkg/logger"
	"go.uber.org/zap"
	"net/http"
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

// CreateTask
//
//	@Summary		CreateTask
//	@Description	Create new task.
//	@UUID			200
//	@Param			params	body		service.CreateTaskParams	true	"Required JSON body with all required task field"
//	@Success		201		{object}	service.CreateTaskResponse	"Task was created successfully"
//	@Failure		400		{object}	response					"Invalid input data"
//	@Failure		500		{object}	response					"Internal error"
//	@Router			/tasks/ [post]
//	@Tags			Task
func (r *taskRoutes) CreateTask(ctx *gin.Context) {
	var params service.CreateTaskParams

	if err := ctx.ShouldBindJSON(&params); err != nil {
		r.logger.Error("error binding json body", zap.Error(err))
		sentErrorResponse(ctx, http.StatusBadRequest, "error binding json body", err)
		return
	}

	resp, err := r.task.CreateTask(ctx, params)
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, constant.ErrInternalError) {
			code = http.StatusInternalServerError
		}
		r.logger.Error("error creating task",
			zap.Error(err),
			zap.String("params", fmt.Sprintf("%+v", params)))
		sentErrorResponse(ctx, code, "error creating task", err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// DeleteTaskByID
//
//	@Summary		DeleteTaskByID
//	@Description	Delete task by its id.
//	@UUID			201
//	@Param			params	path		int			true	"Required task id to delete"
//	@Success		200		{object}	nil			"Task was deleted successfully"
//	@Failure		400		{object}	response	"Invalid input data"
//	@Failure		500		{object}	response	"Internal error"
//	@Router			/tasks/:id [delete]
//	@Tags			Task
func (r *taskRoutes) DeleteTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	err := r.task.DeleteTaskByID(ctx, id)
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, constant.ErrInternalError) {
			code = http.StatusInternalServerError
		}
		r.logger.Error("error deleting task with id",
			zap.Error(err),
			zap.String("task id", id))
		sentErrorResponse(ctx, code, "error deleting task by id", err)
		return
	}

	ctx.Status(http.StatusOK)
	return
}

func (r *taskRoutes) UpdateTaskByID(ctx *gin.Context) {

}

func (r *taskRoutes) GetTaskByID(ctx *gin.Context) {

}
