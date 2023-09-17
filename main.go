package main

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/milosrs/go-hls-server/api"
	"github.com/milosrs/go-hls-server/feed/websocket"
)

const commonEntry = "frontend/dist"

func main() {
	wsHub := websocket.NewHub()
	go wsHub.Start()

	router := api.CreateRouter()
	router.Static("/frontend/dist", "./frontend/dist")
	router.Static("/frontend/src", "./frontend/src")

	router.GET("/", func(ctx *gin.Context) {
		tmp := template.Must(template.ParseFiles(commonEntry + "/index.html"))
		tmp.Execute(ctx.Writer, nil)
	})

	wsRouter := router.Group("/ws")
	wsRouter.GET("", func(ctx *gin.Context) {
		websocket.ServeWS(
			wsHub,
			ctx.Param("topic"),
			ctx.Writer,
			ctx.Request,
		)
		ctx.Status(200)
	})

	router.Run("localhost:8000")

	// Stoping agents
	wsHub.Stop()
}
