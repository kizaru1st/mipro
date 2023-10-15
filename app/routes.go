package app

import (
	"github.com/gorilla/mux"
	"github.com/kizaru1st/mipro/app/controllers"
)

func (server *Server) InitRoutes() {
	server.Router = mux.NewRouter()
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}
