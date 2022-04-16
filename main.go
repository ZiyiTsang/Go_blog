package main

import (
	"Go_blog/bootstrap"
	"Go_blog/pkg/DBTool"
	"Go_blog/pkg/logTool"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
)

var router *mux.Router
var db *sql.DB

func getVariebleFromURL(variable string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[variable]
}

type ArticlesFormData struct {
	Title  string
	Body   string
	URL    *url.URL
	Errors interface{}
	Time   string
	Id     int64
}
type ArticlesData struct {
	Title string
	Body  string
	Time  string //store_time
	Id    int64
}

//func (a ArticlesData) Link() (URL string) {
//	u, err := router.Get("article.show").URL("id", strconv.Itoa(int(a.Id)))
//	logTool.CheckError(err)
//	return u.String()
//}

func (a ArticlesData) delete() (rowaffect int64, err error) {
	deleteSem := "delete from articles where id=?"
	exec, err := db.Exec(deleteSem, a.Id)
	if err != nil {
		return 0, err
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func getArticleByID(id string) (ArticlesData, error) {
	query := "select * from articles where id=?"
	article := ArticlesData{}
	err := db.QueryRow(query, id).Scan(&article.Id, &article.Title, &article.Body, &article.Time)
	return article, err
}
func validateArticleFormData(title string, body string) map[string]string {
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

func saveArticleToDB(title string, body string) (int64, error) {
	var (
		id        int64
		err       error
		result    sql.Result
		statement *sql.Stmt
	)
	statement, err = db.Prepare("insert into articles(title, body,time) VALUES (?,?,now())")
	if err != nil {
		return 0, err
	}
	result, err = statement.Exec(title, body)
	if err != nil {
		return 0, err
	}
	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if id <= 0 {
		return 0, errors.New("id<=0")
	}
	return id, nil
}
func handlerfuncArticlesStore(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "something err in post")
	}
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	errorTag := validateArticleFormData(title, body)
	if len(errorTag) == 0 {
		_, err := fmt.Fprintln(w, "Correct Input data!")
		if err != nil {
			logTool.CheckError(err)
		}
		var increateId int64
		increateId, err = saveArticleToDB(title, body)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, "SQL error!")
			logTool.CheckError(err)
		}
		fmt.Fprintf(w, "insert seccuess full!id=%d\n", increateId)
	} else {
		storeURL, _ := router.Get("articles.store").URL()
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

func handlerfuncArticlesCreate(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()
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
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func handlerfuncArticlesEdit(w http.ResponseWriter, r *http.Request) {
	//  URL:/articles/{id:[0-9]+}/edit
	id := getVariebleFromURL("id", r)
	article, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			fmt.Fprintln(w, "No article")
		} else if err == sql.ErrConnDone {
			w.WriteHeader(500)
			fmt.Fprintln(w, "SQL connection err")
		} else {
			w.WriteHeader(500)
			fmt.Fprintln(w, "other DB fail")
		}
	} else {
		updateURL, _ := router.Get("articles.update").URL("id", id)
		err_tag := make(map[string]string)
		err_tag["title"] = ""
		err_tag["body"] = ""
		data := ArticlesFormData{Title: article.Title, Body: article.Body, URL: updateURL, Time: article.Time, Errors: err_tag, Id: article.Id}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logTool.CheckError(err)
		tmpl.Execute(w, data)
		logTool.CheckError(err)
	}
}

func handlerfuncArticlesUpdate(w http.ResponseWriter, r *http.Request) {
	id := getVariebleFromURL("id", r)
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errorTag := validateArticleFormData(title, body)
	if len(errorTag) == 0 {
		query := "update articles set title=?,body=? where id=?"
		exec, err := db.Exec(query, title, body, id)
		//fmt.Println("1")
		if err != nil {
			//fmt.Println("2")
			w.WriteHeader(500)
			fmt.Fprintln(w, "DB failure in update")
			logTool.CheckError(err)
			return
		} else {
			//fmt.Println("3")
			rowAff, _ := exec.RowsAffected()
			switch rowAff {
			case 0:
				fmt.Fprintln(w, "No any change")

			case 1:
				fmt.Fprintln(w, "change successful")

			}
		}
	} else {
		updateURL, _ := router.Get("articles.update").URL("id", id)
		idNum, _ := strconv.Atoi(id)
		data := ArticlesFormData{Title: title, Body: body, URL: updateURL, Time: "", Id: int64(idNum), Errors: errorTag}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logTool.CheckError(err)
		err = tmpl.Execute(w, data)
		logTool.CheckError(err)
	}
}
func handlerfuncArticlesDelete(w http.ResponseWriter, r *http.Request) {
	id := getVariebleFromURL("id", r)
	article, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(500)
			fmt.Fprint(w, "no such article")
		} else if err == sql.ErrConnDone {
			w.WriteHeader(500)
			fmt.Fprint(w, "SQL connection done")
		} else {
			logTool.CheckError(err)
			w.WriteHeader(500)
			fmt.Fprint(w, "unsolved problem")
		}
	} else {
		rowaff, err := article.delete()
		if err != nil {
			logTool.CheckError(err)
		}
		switch rowaff {
		case 0:
			w.WriteHeader(500)
			fmt.Fprintln(w, "SQL no effect,should no happen")
		case 1:
			fmt.Fprintln(w, "Successful!")

		}
	}
}

func HtmlMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}
func removeTrailingslash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			switch err.(type) {
			//if is runtime error....
			case runtime.Error:
				fmt.Println("runtime:error", err)
			//other error....
			default:
				fmt.Println("other error", err)
			}
			os.Exit(-1)
		}
		fmt.Println("Thank you for using!")
		os.Exit(0)
	}()

	DBTool.Initialize()
	db = DBTool.DB
	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("DB can not close")
			os.Exit(-1)
		}
	}(db)

	fmt.Println("create handle function")

	router.HandleFunc("/articles", handlerfuncArticlesStore).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", handlerfuncArticlesCreate).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", handlerfuncArticlesEdit).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", handlerfuncArticlesUpdate).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", handlerfuncArticlesDelete).Methods("POST").Name("articles.delete")

	router.Use(HtmlMiddleware)
	//start server
	fmt.Println("start server and listening")
	err := http.ListenAndServe(":3000", removeTrailingslash(router))
	if err != nil {
		panic(err)
	}

}
