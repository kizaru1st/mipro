package app

import "github.com/kizaru1st/mipro/app/controllers"

func (server *Server) InitRoutes() {
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}
