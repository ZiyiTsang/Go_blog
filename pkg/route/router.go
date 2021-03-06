package route

import (
	"Go_blog/pkg/logTool"
	"github.com/gorilla/mux"
	"net/http"
)

var route *mux.Router

func SetRoute(r *mux.Router) {
	route = r
}
func Name2URL(routeName string, pairs ...string) string {

	url, err := route.Get(routeName).URL(pairs...)

	if err != nil {
		logTool.CheckError(err)
		return ""
	}
	return url.String()
}
func GetVariebleFromURL(variable string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[variable]
}
