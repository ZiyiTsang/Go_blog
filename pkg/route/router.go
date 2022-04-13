package route

import "github.com/gorilla/mux"

var Router *mux.Router

func Initialize() {
	Router = mux.NewRouter()
}
func RouteName2URL(routerName string, para ...string) string {
	url, err := Router.Get(routerName).URL(para...)
	if err != nil {
		//checkError(err)
		return ""
	}
	return url.String()
}
