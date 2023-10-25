package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/kizaru1st/mipro/app/models"
)

func (server *Server) Products(w http.ResponseWriter, r *http.Request) {
	productModel := models.Product{}
	products, err := productModel.GetProducts(server.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Mengubah objek products menjadi JSON
	jsonResponse, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set header Content-Type menjadi application/json
	w.Header().Set("Content-Type", "application/json")

	// Mengirimkan response JSON ke client
	w.Write(jsonResponse)
}
