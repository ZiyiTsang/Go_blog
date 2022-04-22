package controllers

import (
	"Go_blog/pkg/view"
	"net/http"
)

type AuthController struct{}

func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, ArticlesFormData{}, "auth.register")
}

func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	//
}
