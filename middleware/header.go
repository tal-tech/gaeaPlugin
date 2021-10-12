package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tal-tech/loggerX/logtrace"

	"github.com/tal-tech/gaeaPlugin/xesgin"
)

func RequestHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		//压测标记
		benchmark := c.GetHeader("Request-Type")
		if benchmark == "performance-testing" {
			c.Set("IS_BENCHMARK", "1")
		}
	}
}

func ResponseHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Server", "gosci/"+hostname)
	}
}

func AppInfoHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		ctx := xesgin.TransferToContext(c)
		logtracemap := logtrace.ExtractTraceNodeFromXesContext(ctx)
		device := c.GetHeader("device")
		if device == "" {
			device = c.GetHeader("systemName")
		}
		version := c.GetHeader("version")
		if version == "" {
			version = c.GetHeader("appVersion")
		}
		logtracemap.Set("x_app_device", "\""+device+"\"")
		logtracemap.Set("x_app_version", "\""+version+"\"")
		c.Set(logtrace.GetMetadataKey(), logtracemap)
	}
}
