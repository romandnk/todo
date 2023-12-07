package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/romandnk/todo/pkg/logger"
	"github.com/romandnk/todo/pkg/utils"
	"time"
)

type MW struct {
	logger logger.Logger
}

func NewMiddlewares(logger logger.Logger) *MW {
	return &MW{
		logger: logger,
	}
}

func (m *MW) Logging() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)

		info := utils.RequestInformation(ctx.Request, duration)

		code := ctx.Writer.Status()

		msg := "HTTP requests"

		m.logger.Info(msg,
			"client ip", info.ClientIP,
			"date", info.Date,
			"method", info.Method,
			"method path", info.Path,
			"HTTP version", info.HTTPVersion,
			"status code", code,
			"processing time", info.Latency,
			"user agent", info.UserAgent,
		)
	}
}
