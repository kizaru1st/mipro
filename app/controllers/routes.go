package controllers

import (
	"github.com/gorilla/mux"
)

func (server *Server) InitRoutes() {
	server.Router = mux.NewRouter()
	server.Router.HandleFunc("/", server.Home).Methods("GET")
	server.Router.HandleFunc("/products", server.Products).Methods("GET")
	server.Router.HandleFunc("/products/{slug}", server.GetProductBySlug).Methods("GET")
	server.Router.HandleFunc("/carts", server.GetCart).Methods("GET")
	server.Router.HandleFunc("/carts", server.AddItemToCart).Methods("POST")
	server.Router.HandleFunc("/carts/update", server.UpdateCart).Methods("POST")
	server.Router.HandleFunc("/carts/delete/{id}", server.RemoveItem).Methods("GET")
}
