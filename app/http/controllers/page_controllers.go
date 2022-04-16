package controllers

import (
	"Go_blog/pkg/logTool"
	"fmt"
	"net/http"
)

type PagesController struct {
}

// Home 首页
func (*PagesController) Home(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "<h1>Hello, this is ZIYI's personal Goblog</h1>")
	if err != nil {
		logTool.CheckError(err)
	}
}

// About 关于我们页面
func (*PagesController) About(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "I am Ziyi Tsang,please contact me at:1034337098@qq.com")
	if err != nil {
		logTool.CheckError(err)
	}
}

// NotFound 404 页面
func (*PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "404")
	if err != nil {
		logTool.CheckError(err)
	}
}
