package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func handlerfunc_Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Hello, 这里是 ZZY的goblog</h1>")
}
func handlerfunc_Articles_Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	fmt.Fprint(w, "访问文章列表")

}
func handlerfunc_Articles_Store(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	fmt.Fprint(w, "创建新的文章")

}
func handlerFunc_About(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	fmt.Fprint(w, "I am Ziyi Tsang,please contact me at:1034337098@qq.com")

}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	fmt.Fprint(w, "404")
}
func handlerfunc_Articles_Show(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章 ID："+id)
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", handlerfunc_Root).Methods("Get").Name("home")
	router.HandleFunc("/about", handlerFunc_About).Methods("Get").Name("about")
	//router.HandleFunc("/articles", handlerfunc_Articles_Index).Methods("Get").Name("article.index")
	//router.HandleFunc("/articles", handlerfunc_Articles_Store).Methods("Post").Name("article.store")

	router.HandleFunc("/articles/{id:[0-9]+}", handlerfunc_Articles_Show).Methods("Get").Name("article.show")
	router.HandleFunc("/articles", handlerfunc_Articles_Index).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", handlerfunc_Articles_Store).Methods("POST").Name("articles.store")
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	http.ListenAndServe(":3000", router)
}
