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
	"path/filepath"
	"strconv"
	"time"
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
		viewDir := "resources/views"

		// 2.1 所有布局模板文件 Slice
		files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
		logTool.CheckError(err)

		// 2.2 在 Slice 里新增我们的目标文件
		newFiles := append(files, viewDir+"/articles/index.gohtml")

		// 2.3 解析模板文件
		tmpl, err := template.ParseFiles(newFiles...)
		logTool.CheckError(err)

		// 2.4 渲染模板，将所有文章的数据传输进去
		err = tmpl.ExecuteTemplate(w, "app", articles)
		logTool.CheckError(err)
	}
}

type ArticlesFormData struct {
	Title  string
	Body   string
	URL    string
	Errors map[string]string
	Time   string
	ID     int64
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
		_, err = fmt.Fprint(w, "insert seccuessful\n!")
		logTool.CheckError(err)
		if err != nil {
			return
		}
		_, err = fmt.Fprint(w, article.ID)
		logTool.CheckError(err)
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
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetVariebleFromURL("id", r)
	article, err := article_pkg.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(404)
			_, err := fmt.Fprintln(w, "No article")
			if err != nil {
				logTool.CheckError(err)
			}
		} else {
			w.WriteHeader(500)
			_, err := fmt.Fprintln(w, "other DB fail")
			if err != nil {
				logTool.CheckError(err)
			}
		}
	} else {
		updateURL := route.Name2URL("articles.update", "id", id)
		err_tag := make(map[string]string)
		err_tag["title"] = ""
		err_tag["body"] = ""
		data := ArticlesFormData{Title: article.Title, Body: article.Body, URL: updateURL, Time: article.Time, Errors: err_tag, ID: article.ID}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logTool.CheckError(err)
		err = tmpl.Execute(w, data)
		if err != nil {
			logTool.CheckError(err)
		}
		logTool.CheckError(err)
	}
}
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {

	id := route.GetVariebleFromURL("id", r)
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	errorTag := ValidateArticleFormData(title, body)
	article, _ := article_pkg.Get(id)
	if len(errorTag) == 0 {
		article.Time = time.Now().String()[0:19]
		article.Title = title
		article.Body = body
		rowsAffected, err := article.Update()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprint(w, "update fail")
			if err != nil {
				logTool.CheckError(err)
			}
		} else {
			if rowsAffected == 0 {
				_, err := fmt.Fprint(w, "rowsAffected 0")
				if err != nil {
					logTool.CheckError(err)
				}
			}
			_, err := fmt.Fprint(w, "success!")
			if err != nil {
				logTool.CheckError(err)
			}
		}
	} else {

		updateURL := route.Name2URL("articles.update", "id", id)

		idNum, _ := strconv.Atoi(id)
		data := ArticlesFormData{Title: title, Body: body, URL: updateURL, Time: "", ID: int64(idNum), Errors: errorTag}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logTool.CheckError(err)
		err = tmpl.Execute(w, data)
		logTool.CheckError(err)
	}
}

// Delete 删除文章
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetVariebleFromURL("id", r)
	article, err := article_pkg.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(404)
			_, err := fmt.Fprint(w, "no such article")
			if err != nil {
				return
			}
		} else {
			logTool.CheckError(err)
			w.WriteHeader(500)
			_, err := fmt.Fprint(w, "unsolved problem")
			if err != nil {
				return
			}
		}
	} else {
		roweff, err := article.Delete()
		if err != nil {
			logTool.CheckError(err)
		}
		switch roweff {
		case 0:
			w.WriteHeader(500)
			_, err := fmt.Fprintln(w, "SQL no effect,should no happen")
			if err != nil {
				return
			}
		case 1:
			_, err := fmt.Fprintln(w, "Successful!")
			if err != nil {
				return
			}

		}
	}
}
