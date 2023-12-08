package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romandnk/todo/internal/constant"
	"github.com/romandnk/todo/internal/service"
	statusservice "github.com/romandnk/todo/internal/service/status"
	"github.com/romandnk/todo/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

type statusRoutes struct {
	status service.Status
	logger logger.Logger
}

func newStatusRoutes(g *gin.RouterGroup, status service.Status, logger logger.Logger) {
	r := &statusRoutes{
		status: status,
		logger: logger,
	}

	g.POST("/", r.CreateStatus)
}

// CreateStatus
//
//	@Summary		CreateStatus
//	@Description	Create new task status.
//	@UUID			100
//	@Param			params	body		statusservice.CreateStatusParams	true	"Required JSON body with status name"
//	@Success		200		{object}	statusservice.CreateStatusResponse	"Status was created successfully"
//	@Failure		400		{object}	response							"Invalid input data"
//	@Failure		500		{object}	response							"Internal error"
//	@Router			/statuses/ [post]
//	@Tags			Status
func (r *statusRoutes) CreateStatus(ctx *gin.Context) {
	var params statusservice.CreateStatusParams

	if err := ctx.ShouldBindJSON(&params); err != nil {
		r.logger.Error("error binding json body", zap.Error(err))
		sentErrorResponse(ctx, http.StatusBadRequest, "error binding json body", err)
		return
	}

	resp, err := r.status.CreateStatus(ctx, params)
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, constant.ErrInternalError) {
			code = http.StatusInternalServerError
		}
		r.logger.Error("error creating status", zap.Error(err))
		sentErrorResponse(ctx, code, "error creating status", err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}
