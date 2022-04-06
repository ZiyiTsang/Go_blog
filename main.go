package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var router = mux.NewRouter()

func handlerfunc_Root(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "<h1>Hello, 这里是 ZZY的goblog</h1>")
}
func handlerfunc_Articles_Index(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "访问文章列表")

}
func handlerfunc_Articles_Store(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "创建新的文章")
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
	fmt.Fprint(w, "文章 ID："+id)
}
func handlerfunc_Articles_Create(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>创建文章 —— 我的技术博客</title>
</head>
<body>
    <form action="%s" method="post">
        <p><input type="text" name="title"></p>
        <p><textarea name="body" cols="30" rows="10"></textarea></p>
        <p><button type="submit">提交</button></p>
    </form>
</body>
</html>
`
	storeURL, _ := router.Get("articles.store").URL()
	fmt.Fprintf(w, html, storeURL)
	fmt.Fprint(w, "创建博文表单")
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
		case runtime.Error:
			fmt.Println("runtime:error", err)
		default:
			fmt.Println("other error", err)
		}
		os.Exit(0)
	}()
	router.HandleFunc("/", handlerfunc_Root).Methods("Get").Name("home")
	router.HandleFunc("/about", handlerFunc_About).Methods("Get").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", handlerfunc_Articles_Show).Methods("Get").Name("article.show")
	router.HandleFunc("/articles", handlerfunc_Articles_Index).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", handlerfunc_Articles_Store).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", handlerfunc_Articles_Create).Methods("GET").Name("articles.create")
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.Use(HTML_Middleware)
	err := http.ListenAndServe(":3000", remove_TrailingSlash(router))
	if err != nil {
		panic(err)
	}
}
