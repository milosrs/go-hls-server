package html_handlers

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func RenderHome(ctx *gin.Context) {
	tmpl := template.Must(template.ParseFiles("frontend/index.html"))
	ctx.Writer.Header().Add("Cache-Control", "no-cache, private, max-age=0")
	ctx.Writer.Header().Add("Pragma", "no-cache")
	ctx.Writer.Header().Add("X-Accel-Expires", "0")
	tmpl.Execute(ctx.Writer, nil)
}
