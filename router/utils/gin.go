package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/zzzgydi/zbyai/common"
	"github.com/zzzgydi/zbyai/common/logger"
)

func GetTraceLogger(c *gin.Context) *logger.TraceLogger {
	if trace, ok := c.Get(common.CTX_TRACE_LOGGER); ok {
		if trace, ok := trace.(*logger.TraceLogger); ok {
			return trace
		}
	}
	return nil
}

func GetRealUserIp(c *gin.Context) string {
	if ip, ok := c.Get(common.CTX_REAL_IP); ok {
		if ip, ok := ip.(string); ok {
			return ip
		}
	}
	ip := c.ClientIP()
	c.Set(common.CTX_REAL_IP, ip)
	return ip
}

func realIp(c *gin.Context) string {
	r := c.Request
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("Fly-Client-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	// parse X-Forwarded-For and trust platform
	return c.ClientIP()
}
