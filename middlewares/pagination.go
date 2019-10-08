package middlewares

import (
	"github.com/kataras/iris"
)

func Paginate(ctx iris.Context) {
	page, err := ctx.URLParamInt("page")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": "Invalid page provided",
		})
		return
	}

	limit, err := ctx.URLParamInt("limit")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": "Invalid limit provided",
		})
		return
	}

	if page < 1 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": "limit should be greater or equal to 1",
		})
		return
	}

	if limit < 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": "limit should be greater or equal to 0",
		})
		return
	}

	offset := (page - 1) * limit

	ctx.Values().SetImmutable("_page", page)
	ctx.Values().SetImmutable("_limit", limit)
	ctx.Values().SetImmutable("_offset", offset)

	ctx.Next()
}
