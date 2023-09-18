package files

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/milosrs/go-hls-server/files/model"
)

// 30mb max request
const maxRequest = 30000000

func setCookiesToResponse(name string, length float64, writer http.ResponseWriter) {
	fileNameCookie := &http.Cookie{
		Name:     "file-name",
		Value:    name,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	lengthCookie := &http.Cookie{
		Name:     "uploaded-length",
		Value:    fmt.Sprintf("%d", length),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(writer, fileNameCookie)
	http.SetCookie(writer, lengthCookie)
}

func ReceiveFirstChunk(s Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fileData := model.InitialFileData{}
		err := ctx.ShouldBind(&fileData)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		length, err := s.CreateFile(&fileData)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		setCookiesToResponse(fileData.Name, length, ctx.Writer)
		ctx.String(http.StatusOK, "%d", length)
	}
}

func PatchChunk(s Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data model.FileChunk
		err := ctx.ShouldBind(&data)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
		}

		length, err := s.AppendChunk(&data)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		setCookiesToResponse(data.Name, length, ctx.Writer)
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
