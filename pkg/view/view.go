package view

import (
	"Go_blog/pkg/logTool"
	"Go_blog/pkg/route"
	"Go_blog/pkg/typesTool"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

func Render(w io.Writer, routerName string, data interface{}) { //routerName eg."articles.show"
	viewDir := "resources/views"
	files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
	logTool.CheckError(err)
	name := strings.Replace(routerName, ".", "/", -1)
	name += ".gohtml"
	name = "/" + name
	newFiles := append(files, viewDir+name)
	tmpl, err := template.New(name).
		Funcs(template.FuncMap{
			"RouteName2URL":  route.Name2URL,
			"Uint64ToString": typesTool.Int64ToString,
		}).ParseFiles(newFiles...)
	logTool.CheckError(err)
	err = tmpl.ExecuteTemplate(w, "app", data)
	logTool.CheckError(err)
}
