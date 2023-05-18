package api

import (
	"github.com/gin-gonic/gin"
	"github.com/milosrs/go-hls-server/files"
)

const (
	apiGroup    = "/api"
	usersGroup  = "/users"
	uploadChunk = "/uploadChunk"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group(apiGroup)
	users := api.Group(usersGroup)

	fileService := files.ServiceImpl{}

	users.PATCH(uploadChunk, files.PatchChunk(&fileService))
	users.DELETE(uploadChunk, files.RemoveChunk(&fileService))

	return router
}
