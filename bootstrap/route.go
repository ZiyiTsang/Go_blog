package bootstrap

import (
	"Go_blog/pkg/route"
	"Go_blog/routes"
	"fmt"
	"github.com/gorilla/mux"
)

func SetupRoute() *mux.Router {
	fmt.Println("initiate route...")
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)
	route.SetRoute(router)
	return router
}
