package main

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/milosrs/go-hls-server/api"
	"github.com/milosrs/go-hls-server/feed/websocket"
	ws "golang.org/x/net/websocket"
)

const commonEntry = "frontend/dist"

func main() {
	wsServer := websocket.NewServer()

	router := api.CreateRouter()
	router.Static("/frontend/dist", "./frontend/dist")
	router.Static("/frontend/src", "./frontend/src")

	router.GET("/", func(ctx *gin.Context) {
		tmp := template.Must(template.ParseFiles(commonEntry + "/index.html"))
		tmp.Execute(ctx.Writer, nil)
	})
	router.GET("/ws", func(ctx *gin.Context) {
		result := wsServer.HandleUpgrade(*ctx.Request)
		ctx.Request.Response = result
		ctx.Status(200)
	})

	mux := http.NewServeMux()
	mux.Handle("/", ws.Handler(wsServer.HandleConnection))

	router.Run()
}
