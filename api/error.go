package api

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func successOrAbort(ctx *gin.Context, code int, err error) (success bool) {
	if err != nil {
		ctx.AbortWithError(code, err)
	}
	return err == nil
}
