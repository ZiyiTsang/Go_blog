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

func Render(w io.Writer, data interface{}, tplFiles ...string) {
	RenderTemplate(w, "app", data, tplFiles...)
}

// RenderSimple 渲染简单的视图
func RenderSimple(w io.Writer, data interface{}, tplFiles ...string) {
	RenderTemplate(w, "simple", data, tplFiles...)
}

func RenderTemplate(w io.Writer, name string, data interface{}, tplFiles ...string) {
	viewDir := "resources/views/"
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}
	files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
	logTool.CheckError(err)
	allFiles := append(files, tplFiles...)
	tmpl, err := template.New(name).
		Funcs(template.FuncMap{
			"RouteName2URL":  route.Name2URL,
			"Uint64ToString": typesTool.Int64ToString,
		}).ParseFiles(allFiles...)
	logTool.CheckError(err)
	err = tmpl.ExecuteTemplate(w, name, data)
	logTool.CheckError(err)
}
