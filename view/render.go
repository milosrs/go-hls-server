package view

import (
	"bytes"
	"html/template"
	"log"
)

const componentsPath = "frontend/src/components/"

func createComponentPath(name component) (path string) {
	path = componentsPath +
		string(name) +
		"/" +
		string(name) +
		".html"

	return
}

func CreateComponentData(name component, data any) string {
	buf := bytes.NewBufferString("")
	path := createComponentPath(name)

	tmpl, err := template.New(string(name)).Parse(path)
	if err != nil {
		log.Printf("error loading file: %v", err)
		return ""
	}

	err = tmpl.Execute(buf, data)
	if err != nil {
		log.Printf("error executing template file: %v", err)
		return ""
	}

	return buf.String()
}
