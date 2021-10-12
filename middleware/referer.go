package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tal-tech/loggerX"
)

const (
	MisMatchStop = iota
	MisMatchNext
)

// CheckReferer
// @param url 需校验的url
// @param isStop 不匹配时是否是否继续执行,MisMatchNext:是,MisMatchStop:否
// return gin.HandlerFunc
func CheckReferer(url string, isStop int) gin.HandlerFunc {
	logTag := "http.middleware.refer"
	return func(ctx *gin.Context) {
		if ref := ctx.GetHeader("Referer"); strings.Contains(ref, url) {
			ctx.Next()
			return
		}
		switch isStop {

		case MisMatchNext:
			logger.Wx(ctx, logTag, "refer is Illegal Next url is %s, refer is :%s", ctx.FullPath(), ctx.GetHeader("Referer"))
			ctx.Next()
			return
		case MisMatchStop:
			logger.Ex(ctx, logTag, "refer is Illegal Stop url is %s, refer is :%s", ctx.FullPath(), ctx.GetHeader("Referer"))
			ctx.AbortWithStatus(403)
			return
		}
	}
}
