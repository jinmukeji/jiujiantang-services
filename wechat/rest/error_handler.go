package rest

import (
	"github.com/kataras/iris/v12"
)

func internalServerError(ctx iris.Context) {
	// nolint: errcheck, gas
	ctx.JSON(iris.Map{
		"status":      iris.StatusInternalServerError, // 500
		"description": "Internal Server Error",
	})
}


func writeSessionErrorJSON(ctx iris.Context, redirectURL string, err error) {
	// nolint: errcheck, gas
	ctx.JSON(iris.Map{
        "ok":       false,
		"error":    wrapError(ErrSessionUnauthorized,"",err),
		"redirect": redirectURL,
	})
}
