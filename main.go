package main

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/milosrs/go-hls-server/api"
	"github.com/milosrs/go-hls-server/feed/websocket"
	"github.com/milosrs/go-hls-server/files"
)

const commonEntry = "frontend/dist"

func main() {
	hub := websocket.NewHub()
	fileService := files.NewService()

	go hub.Start()
	go fileService.Start()

	router := api.CreateRouter()
	router.Static("/frontend/dist", "./frontend/dist")
	router.Static("/frontend/src", "./frontend/src")

	router.GET("/", func(ctx *gin.Context) {
		tmp := template.Must(template.ParseFiles(commonEntry + "/index.html"))
		tmp.Execute(ctx.Writer, nil)
	})

	wsRouter := router.Group("/ws")
	wsRouter.GET("/:topic", func(ctx *gin.Context) {
		conn := websocket.ServeWS(ctx.Writer, ctx.Request)
		files.NewFeed(conn, hub, fileService)
		ctx.Status(200)
	})

	router.Run("localhost:8000")

	// Stoping agents
	hub.Stop()
}
