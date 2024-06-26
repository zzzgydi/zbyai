package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zzzgydi/zbyai/common"
	"github.com/zzzgydi/zbyai/common/logger"
	"github.com/zzzgydi/zbyai/router/utils"
)

func LoggerMiddleware(c *gin.Context) {
	trace := logger.NewTraceLogger(c)
	trace.SetIp(utils.GetRealUserIp(c))
	c.Set(common.CTX_TRACE_LOGGER, trace)
	c.Header("X-Trace-Id", trace.RequestId)
	c.Next()
	trace.Write()
}
