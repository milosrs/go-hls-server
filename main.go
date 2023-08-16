package main

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/milosrs/go-hls-server/api"
)

const commonEntry = "frontend/dist"

func main() {
	router := api.CreateRouter()

	router.Static("/frontend/dist", "./frontend/dist")
	router.Static("/frontend/src", "./frontend/src")

	router.GET("/", func(ctx *gin.Context) {
		tmp := template.Must(template.ParseFiles(commonEntry + "/index.html"))
		tmp.Execute(ctx.Writer, nil)
	})

	router.Run()
}
