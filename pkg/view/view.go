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

func Render(w io.Writer, routerName string, data interface{}, tplFiles ...string) { //routerName eg."articles.show"->extract the relate template files
	viewDir := "resources/views/"
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}
	files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
	logTool.CheckError(err)
	name := strings.Replace(routerName, ".", "/", -1)
	name += ".gohtml"
	rudimentaryFiles := append(files, viewDir+name)
	allFiles := append(rudimentaryFiles, tplFiles...)
	tmpl, err := template.New(name).
		Funcs(template.FuncMap{
			"RouteName2URL":  route.Name2URL,
			"Uint64ToString": typesTool.Int64ToString,
		}).ParseFiles(allFiles...)
	logTool.CheckError(err)
	err = tmpl.ExecuteTemplate(w, "app", data)
	logTool.CheckError(err)
}
