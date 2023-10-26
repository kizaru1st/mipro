package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/kizaru1st/mipro/app/models"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

func GetShoppingCartID(w http.ResponseWriter, r *http.Request) string {
	session, _ := store.Get(r, sessionShoppingCart)
	if session.Values["cart-id"] == nil {
		cartID := uuid.New().String()
		session.Values["cart-id"] = cartID
		_ = session.Save(r, w)
		return cartID
	}

	return session.Values["cart-id"].(string)
}

func GetShoppingCart(db *gorm.DB, cartID string) (*models.Cart, error) {
	var cart models.Cart

	existCart, err := cart.GetCart(db, cartID)
	if err != nil {
		existCart, _ = cart.CreateCart(db, cartID)
	}

	return existCart, nil
}

func (server *Server) GetCart(w http.ResponseWriter, r *http.Request) {
	cartID := GetShoppingCartID(w, r)
	cart, err := GetShoppingCart(server.DB, cartID)

	if err != nil {
		jsonResponse := map[string]string{"message": "Error getting the cart"}
		sendJSONResponse(w, jsonResponse, http.StatusInternalServerError)
		return
	}

	cartJSON, err := json.Marshal(cart)
	if err != nil {
		jsonResponse := map[string]string{"message": "Error converting cart to JSON"}
		sendJSONResponse(w, jsonResponse, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(cartJSON)
	if err != nil {
		jsonResponse := map[string]string{"message": "Error sending JSON response"}
		sendJSONResponse(w, jsonResponse, http.StatusInternalServerError)
		return
	}
}

func (server *Server) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	productID := r.FormValue("product_id")
	qty, _ := strconv.Atoi(r.FormValue("qty"))

	productModel := models.Product{}
	product, err := productModel.FindByID(server.DB, productID)
	if err != nil {
		jsonResponse := map[string]string{"message": "Product not found"}
		sendJSONResponse(w, jsonResponse, http.StatusNotFound)
		return
	}

	if qty > product.Stock {
		jsonResponse := map[string]string{"message": "Quantity exceeds available stock"}
		sendJSONResponse(w, jsonResponse, http.StatusBadRequest)
		return
	}

	cartID := GetShoppingCartID(w, r)
	cart, err := GetShoppingCart(server.DB, cartID)
	if err != nil {
		jsonResponse := map[string]string{"message": "Error getting the cart"}
		sendJSONResponse(w, jsonResponse, http.StatusInternalServerError)
		return
	}

	//Tambahkan barang ke keranjang
	fmt.Println("cart id ===> ", cart.ID)
	jsonResponse := map[string]string{"message": "Product added to cart"}
	sendJSONResponse(w, jsonResponse, http.StatusOK)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonData)
}

func (server *Server) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	productModel := models.Product{}
	product, err := productModel.FindByID(server.DB, productID)
	if err != nil {
		// Handle error
		jsonResponse := map[string]string{"message": "Product not found"}
		sendJSONResponse(w, jsonResponse, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(product); err != nil {
		// Handle JSON encoding error
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
