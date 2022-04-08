package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" //Anonymous import->enable support for MySQL,,but not use directly
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"unicode/utf8"
)

var router = mux.NewRouter()

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
	if len(error_tag) == 0 {
		fmt.Fprint(w, "Correct posting!")
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
func main() {
	defer func() {
		err := recover()
		switch err.(type) {
		//if is runtime error....
		case runtime.Error:
			fmt.Println("runtime:error", err)
		//other error....
		default:
			fmt.Println("other error", err)
		}
		os.Exit(0)
	}()
	//create relation between address and handle_function
	router.HandleFunc("/", handlerfunc_Root).Methods("Get").Name("home")
	router.HandleFunc("/about", handlerFunc_About).Methods("Get").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", handlerfunc_Articles_Show).Methods("Get").Name("article.show")
	router.HandleFunc("/articles", handlerfunc_Articles_Index).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", handlerfunc_Articles_Store).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", handlerfunc_Articles_Create).Methods("GET").Name("articles.create")
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.Use(HTML_Middleware)
	//start server
	err := http.ListenAndServe(":3000", remove_TrailingSlash(router))
	if err != nil {
		panic(err)
	}
}
