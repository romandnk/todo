package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/romandnk/todo/internal/service"
	"github.com/romandnk/todo/pkg/logger"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	engine   *gin.Engine
	services *service.Services
	logger   logger.Logger
	mw       *MW
}

func NewHandler(services *service.Services, logger logger.Logger, mw *MW) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
		mw:       mw,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	h.engine = router

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := router.Group("/api/v1", h.mw.Logging())
	{
		// status management group
		statuses := api.Group("/statuses")
		{
			newStatusRoutes(statuses, h.services.Status, h.logger)
		}

		// task management group
	}

	return router
}
