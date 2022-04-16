package bootstrap

import (
	"Go_blog/pkg/route"
	"Go_blog/routes"
	"github.com/gorilla/mux"
)

func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)
	route.SetRoute(router)
	return router
}
