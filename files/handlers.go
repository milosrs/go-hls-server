package files

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/milosrs/go-hls-server/files/model"
)

// 30mb max request
const maxRequest = 30000000
const fileName = "name"

func PatchChunk(s Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data *model.File
		err := ctx.ShouldBind(data)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
		}

		length, err := s.AppendChunk(data)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		ctx.JSON(http.StatusOK, length)
	}
}

func RemoveChunk(s Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("fileName")

		err := s.Remove(name)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		ctx.Status(http.StatusOK)
	}
}
