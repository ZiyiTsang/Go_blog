package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 这里是 ZZY的goblog</h1>")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "I am ZZY,please contact me at:1034337098@qq.com")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "wrong!!Please try again")
	}
}
func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
