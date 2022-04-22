package controllers

import (
	articlepkg "Go_blog/app/models/article"
	"Go_blog/pkg/logTool"
	"Go_blog/pkg/route"
	"Go_blog/pkg/view"
	"fmt"
	"gorm.io/gorm"
	"net/http"
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
	article, err := articlepkg.Get(id)

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
		view.Render(w, article, "articles.show")
	}
}

// Index 文章列表页
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := articlepkg.GetAll()
	if err != nil {
		logTool.CheckError(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprint(w, "500")
		if err != nil {
			logTool.CheckError(err)
		}
	} else {
		view.Render(w, articles, "articles.index")
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
	view.Render(w, data, "articles.create")
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
		article := articlepkg.Article{
			Title: title,
			Body:  body,
		}
		var rowEff int64
		rowEff, err = article.CreateWithTitleBody()
		if err != nil {
			w.WriteHeader(500)
			_, err := fmt.Fprint(w, "SQL error!")
			if err != nil {
				logTool.CheckError(err)
			}
			logTool.CheckError(err)
		}
		if rowEff <= 0 {
			fmt.Fprintf(w, "row effected fail")
		}
		_, err = fmt.Fprint(w, "insert successful\n!")
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
		view.Render(w, data, "articles.store", "articles._form_field")
	}
}
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetVariebleFromURL("id", r)
	article, err := articlepkg.Get(id)
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
		errTag := make(map[string]string)
		errTag["title"] = ""
		errTag["body"] = ""
		data := ArticlesFormData{Title: article.Title, Body: article.Body, URL: updateURL, Time: article.Time, Errors: errTag, ID: article.ID}
		view.Render(w, data, "articles.edit", "articles._form_field")
	}
}
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {

	id := route.GetVariebleFromURL("id", r)
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	errorTag := ValidateArticleFormData(title, body)
	article, _ := articlepkg.Get(id)
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
		view.Render(w, data, "articles.edit", "articles._form_field")
	}
}

// Delete 删除文章
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetVariebleFromURL("id", r)
	article, err := articlepkg.Get(id)
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
