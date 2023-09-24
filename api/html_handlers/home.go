package html_handlers

import (
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/milosrs/go-hls-server/files/model"
	"github.com/milosrs/go-hls-server/view"
	"github.com/milosrs/go-hls-server/view/alert"
	"github.com/milosrs/go-hls-server/view/homepage"
)

func execHome(ctx *gin.Context, data homepage.Def) {
	tmpl := template.Must(template.ParseFiles("frontend/index.html"))
	ctx.Writer.Header().Add("Cache-Control", "no-cache, private, max-age=0")
	ctx.Writer.Header().Add("Pragma", "no-cache")
	ctx.Writer.Header().Add("X-Accel-Expires", "0")
	tmpl.Execute(ctx.Writer, data)
}

func RenderHome(ctx *gin.Context) {
	home := homepage.Def{
		FileUpload:    true,
		Progress:      false,
		ProgressSteps: false,
		Alerts:        []string{},
	}

	execHome(ctx, home)
}

func FileHandshake(ctx *gin.Context) {
	defer ctx.Request.Body.Close()

	fd := model.FileDescription{}

	err := ctx.ShouldBind(&fd)
	if err != nil {
		log.Printf("Error in binding file", err)
		alert := alert.Alert{
			Color:       view.Red,
			Heading:     "Error processing request!",
			Description: "Your data has a wrong format.",
		}
		data := view.CreateComponentData("alert", alert)
		home := homepage.Def{
			FileUpload:    true,
			Progress:      false,
			ProgressSteps: false,
			Alerts:        []string{data},
		}

		execHome(ctx, home)
	}
}
