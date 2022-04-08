package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" //Anonymous import->enable support for MySQL,,but not use directly
	"github.com/gorilla/mux"
	"github.com/zalando/go-keyring"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"
)

var router = mux.NewRouter()
var db *sql.DB

func execSql(s string) {
	_, err := db.Exec(s)
	checkError(err)
	fmt.Println("exec order successful")
}
func handlerfunc_Root(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "<h1>Hello, this is ZIYI's personal Goblog</h1>")
}
func handlerfunc_Articles_Index(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "article index")

}

type ArticlesFormData struct {
	Title  string
	Body   string
	URL    *url.URL
	Errors interface{}
}

func save_article_to_db(title string, body string) (int64, error) {
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
	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if id <= 0 {
		return 0, errors.New("id<=0")
	}
	return id, nil
}
func handlerfunc_Articles_Store(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "something err in post")
	}
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	error_tag := make(map[string]string)
	if title == "" {
		error_tag["title"] = "Title:empty"
	} else if utf8.RuneCountInString(title) > 40 || utf8.RuneCountInString(title) < 8 {
		error_tag["title"] = "Title:too long/too short(needs 8-40)"
	}
	if body == "" {
		error_tag["body"] = "Body:empty"
	} else if utf8.RuneCountInString(body) > 200 || utf8.RuneCountInString(body) < 8 {
		error_tag["body"] = "Body:too long/too short(needs 8-20)"
	}
	//
	if len(error_tag) == 0 {
		_, err := fmt.Fprintln(w, "Correct posting!")
		if err != nil {
			checkError(err)
		}
		var increate_id int64
		increate_id, err = save_article_to_db(title, body)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, "SQL error!")
			checkError(err)
		}
		fmt.Fprintf(w, "insert seccuess full!id=%d\n", increate_id)
	} else {
		storeURL, _ := router.Get("articles.store").URL()
		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: error_tag,
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

func handlerFunc_About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "I am Ziyi Tsang,please contact me at:1034337098@qq.com")
}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "404")
}
func handlerfunc_Articles_Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "Article ID："+id)
}
func handlerfunc_Articles_Create(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()
	err_tag := make(map[string]string)
	err_tag["title"] = ""
	err_tag["body"] = ""
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: err_tag,
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

func HTML_Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}
func remove_TrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func initDB() {
	var err error
	mysql_passwd, err := keyring.Get("mysql", "root")
	checkError(err)
	mysql_address, err := keyring.Get("mysql", "address")
	checkError(err)
	config := mysql.Config{
		User:                 "root",
		Passwd:               mysql_passwd,
		Addr:                 mysql_address,
		Net:                  "tcp",
		DBName:               "go_blog",
		AllowNativePasswords: true,
		Timeout:              time.Hour * 2,
		CheckConnLiveness:    true,
	}
	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)
	//my mySQL "wait_timeout" shows "7200"(s)=2hour,I set same as it did..
	db.SetConnMaxLifetime(2 * time.Hour)
	//my mySQL "max_connections" shows 2520,so I set 2000 here..
	db.SetMaxOpenConns(1000)
	//I think it is ok for more than 10..
	db.SetMaxIdleConns(40)
	err = db.Ping()
	checkError(err)
	fmt.Println("init DB successful")
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
		os.Exit(0)
	}()
	initDB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("DB can not close")
			os.Exit(-1)
		}
	}(db)
	//create relation between address and handle_function
	fmt.Println("create handle function")
	router.HandleFunc("/", handlerfunc_Root).Methods("Get").Name("home")
	router.HandleFunc("/about", handlerFunc_About).Methods("Get").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", handlerfunc_Articles_Show).Methods("Get").Name("article.show")
	router.HandleFunc("/articles", handlerfunc_Articles_Index).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", handlerfunc_Articles_Store).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", handlerfunc_Articles_Create).Methods("GET").Name("articles.create")
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.Use(HTML_Middleware)
	//start server
	fmt.Println("start server and listening")
	err := http.ListenAndServe(":3000", remove_TrailingSlash(router))
	if err != nil {
		panic(err)
	}

}
