package controllers

import (
	article_pkg "Go_blog/app/models/article"
	"Go_blog/pkg/logTool"
	"Go_blog/pkg/route"
	"Go_blog/pkg/typesTool"
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

type ArticlesController struct {
}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 1. get URL id
	id := route.GetVariebleFromURL("id", r)

	// 2. load article from mysql
	article, err := article_pkg.Get(id)

	// 3. if wrong
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 article not found")
		} else {
			logTool.CheckError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 server down")
		}
	} else {
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL": route.Name2URL,
				"Int64ToString": typesTool.Int64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		logTool.CheckError(err)
		tmpl.Execute(w, article)
		logTool.CheckError(err)
	}
}
