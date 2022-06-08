package ginredoc

import (
	"github.com/gin-gonic/gin"
	"github.com/mvrilo/go-redoc"
)

func GinHandler(r redoc.Redoc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r.Handler()(ctx.Writer, ctx.Request)
	}
}
