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
	"unicode/utf8"
)

func ValidateArticleFormData(title string, body string) map[string]string {
	e := make(map[string]string)
	if title == "" {
		e["title"] = "title can not be empty"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		e["title"] = "title should within 3-40 character"
	}

	if body == "" {
		e["body"] = "content cannot be empty"
	} else if utf8.RuneCountInString(body) < 10 || utf8.RuneCountInString(body) > 500 {
		e["body"] = "content should within 10-500 character"
	}
	return e
}

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
			_, err := fmt.Fprint(w, "404 article not found")
			if err != nil {
				logTool.CheckError(err)
			}
		} else {
			logTool.CheckError(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprint(w, "500 server down")
			if err != nil {
				logTool.CheckError(err)
			}
		}
	} else {
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL": route.Name2URL,
				"Int64ToString": typesTool.Int64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		logTool.CheckError(err)
		err = tmpl.Execute(w, article)
		if err != nil {
			logTool.CheckError(err)
		}
		logTool.CheckError(err)
	}
}

// Index 文章列表页
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := article_pkg.GetAll()
	if err != nil {
		logTool.CheckError(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprint(w, "500")
		if err != nil {
			logTool.CheckError(err)
		}
	} else {
		tmpl, _ := template.ParseFiles("resources/views/articles/index.gohtml")
		err = tmpl.Execute(w, articles)
		logTool.CheckError(err)
	}
}

type ArticlesFormData struct {
	Title  string
	Body   string
	URL    string
	Errors map[string]string
}

func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	storeURL := route.Name2URL("articles.store")
	errTag := make(map[string]string)
	errTag["title"] = ""
	errTag["body"] = ""
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: errTag,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		logTool.CheckError(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		logTool.CheckError(err)
	}

}

// Store 文章创建页面
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		_, err := fmt.Fprint(w, "something err in post")
		if err != nil {
			logTool.CheckError(err)
		}
	}
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	errorTag := ValidateArticleFormData(title, body)
	if len(errorTag) == 0 {
		_, err := fmt.Fprint(w, "Correct Input data!\n")
		if err != nil {
			logTool.CheckError(err)
		}
		article := article_pkg.Article{
			Title: title,
			Body:  body,
		}
		var row_eff int64
		row_eff, err = article.CreateWithTitleBody()
		if err != nil {
			w.WriteHeader(500)
			_, err := fmt.Fprint(w, "SQL error!")
			if err != nil {
				logTool.CheckError(err)
			}
			logTool.CheckError(err)
		}
		if row_eff <= 0 {
			fmt.Fprintf(w, "row effected fail")
		}
		fmt.Fprint(w, "insert seccuessful\n!")
		fmt.Fprint(w, article.ID)
	} else {
		storeURL := route.Name2URL("articles.store")
		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errorTag,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}
