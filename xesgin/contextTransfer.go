package xesgin

import (
	"context"

	"github.com/gin-gonic/gin"
)

func TransferToContext(c *gin.Context) context.Context {
	ctx := context.Background()
	for k, v := range c.Keys {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}
