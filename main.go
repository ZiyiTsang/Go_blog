package main

import (
	"Go_blog/bootstrap"
	"Go_blog/pkg/logTool"
	"database/sql"
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
	fmt.Println("initiating...")
	//DBTool.Initialize()
	//db = DBTool.DB
	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("DB can not close")
			os.Exit(-1)
		}
	}(db)
	router.Use(HtmlMiddleware)

	go func() {
		err := http.ListenAndServe(":3000", removeTrailingslash(router))
		panic(err)
	}()
	var choose string
	for {
		out := 0
		fmt.Println("Do you want to exit?")
		_, err := fmt.Scan(&choose)
		if err != nil {
			logTool.CheckError(err)
		}
		switch choose {
		case "0":
			out = 1
		default:
			fmt.Println("input again")
		}
		if out == 1 {
			break
		}
	}
}
