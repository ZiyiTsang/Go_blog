package route

import (
	"Go_blog/routes"
	"github.com/gorilla/mux"
	"net/http"
)

var Router *mux.Router

func Initialize() {
	Router = mux.NewRouter()
	routes.RegisterWebRoutes(Router)
}
func RouteName2URL(routerName string, para ...string) string {
	url, err := Router.Get(routerName).URL(para...)
	if err != nil {
		//checkError(err)
		return ""
	}
	return url.String()
}
func GetVariebleFromURL(variable string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[variable]
}
