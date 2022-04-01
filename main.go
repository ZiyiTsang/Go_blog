package main

import (
	"fmt"
	"net/http"
)

func handlerFunc_root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 这里是 ZZY的goblog</h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "wrong!!Please try again(home)")
	}
}
func handlerFunc_articles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	if r.URL.Path == "/articles/" {
		fmt.Fprint(w, string(r.URL.Path))
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "wrong!!Please try again(article)")
	}
}
func handlerFunc_about(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	if r.URL.Path == "/about" {
		fmt.Fprint(w, "I am Ziyi Tsang,please contact me at:1034337098@qq.com")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404!(about)")
	}
}
func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", handlerFunc_root)
	router.HandleFunc("/about", handlerFunc_about)
	router.HandleFunc("/articles/", handlerFunc_articles)
	http.ListenAndServe(":3000", router)
}
