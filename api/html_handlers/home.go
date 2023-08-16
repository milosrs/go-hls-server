package html_handlers

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func RenderHome(ctx *gin.Context) {
	tmpl := template.Must(template.ParseFiles("frontend/index.html"))
	tmpl.Execute(ctx.Writer, nil)
}
