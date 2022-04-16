package main

import (
	"Go_blog/bootstrap"
	"Go_blog/pkg/DBTool"
	"Go_blog/pkg/logTool"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
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

func handlerfuncArticlesDelete(w http.ResponseWriter, r *http.Request) {
	id := getVariebleFromURL("id", r)
	article, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(500)
			_, err := fmt.Fprint(w, "no such article")
			if err != nil {
				return
			}
		} else if err == sql.ErrConnDone {
			w.WriteHeader(500)
			_, err := fmt.Fprint(w, "SQL connection done")
			if err != nil {
				logTool.CheckError(err)
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
		rowaff, err := article.delete()
		if err != nil {
			logTool.CheckError(err)
		}
		switch rowaff {
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

	router.HandleFunc("/articles/{id:[0-9]+}/delete", handlerfuncArticlesDelete).Methods("POST").Name("articles.delete")

	router.Use(HtmlMiddleware)
	//start server
	fmt.Println("start server and listening")
	err := http.ListenAndServe(":3000", removeTrailingslash(router))
	if err != nil {
		panic(err)
	}

}
