package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

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

func (server *Server) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["slug"] == "" {
		http.Error(w, "Slug tidak valid", http.StatusBadRequest)
		return
	}

	productModel := models.Product{}
	product, err := productModel.FindBySlug(server.DB, vars["slug"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode produk menjadi JSON
	jsonData, err := json.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set header respons untuk memberitahu bahwa kita mengirimkan JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Mengirimkan data JSON sebagai respons
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
