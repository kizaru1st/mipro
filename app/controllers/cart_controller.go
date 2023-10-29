package controllers

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/kizaru1st/mipro/app/models"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

func (server *Server) respondWithError(w http.ResponseWriter, code int, message string) {
	server.respondWithJSON(w, code, map[string]string{"error": message})
}

func (server *Server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func GetShoppingCartID(w http.ResponseWriter, r *http.Request) string {
	session, _ := store.Get(r, sessionShoppingCart)
	if session.Values["cart-id"] == nil {
		session.Values["cart-id"] = uuid.New().String()
		_ = session.Save(r, w)
	}

	return fmt.Sprintf("%v", session.Values["cart-id"])
}

func GetShoppingCart(db *gorm.DB, cartID string) (*models.Cart, error) {
	var cart models.Cart

	existCart, err := cart.GetCart(db, cartID)
	if err != nil {
		existCart, err = cart.CreateCart(db, cartID)
		if err != nil {
			return nil, err
		}
	}

	_, err = existCart.CalculateCart(db, cartID)
	if err != nil {
		return nil, err
	}

	updatedCart, err := cart.GetCart(db, cartID)
	if err != nil {
		return nil, err
	}

	return updatedCart, nil
}

func (server *Server) GetCart(w http.ResponseWriter, r *http.Request) {
	var cart *models.Cart

	cartID := GetShoppingCartID(w, r)
	cart, err := GetShoppingCart(server.DB, cartID)
	if err != nil {
		// Handle error, misalnya dengan mengembalikan respons JSON
		server.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	items, err := cart.GetItems(server.DB, cartID)
	if err != nil {
		// Handle error, misalnya dengan mengembalikan respons JSON
		server.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Mengembalikan respons JSON
	server.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"cart":  cart,
		"items": items,
	})
}

func (server *Server) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	productID := r.FormValue("product_id")
	qty, _ := strconv.Atoi(r.FormValue("qty"))

	productModel := models.Product{}
	product, err := productModel.FindByID(server.DB, productID)
	if err != nil {
		// Handle error, misalnya dengan mengembalikan respons JSON
		server.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if qty > product.Stock {
		// Handle error, misalnya dengan mengembalikan respons JSON
		server.respondWithError(w, http.StatusBadRequest, "Quantity exceeds available stock")
		return
	}

	var cart *models.Cart

	cartID := GetShoppingCartID(w, r)
	cart, err = GetShoppingCart(server.DB, cartID)
	if err != nil {
		// Handle error, misalnya dengan mengembalikan respons JSON
		server.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	_, err = cart.AddItem(server.DB, models.CartItem{
		ProductID: productID,
		Qty:       qty,
	})
	if err != nil {
		// Handle error, misalnya dengan mengembalikan respons JSON
		server.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Mengembalikan respons JSON
	server.respondWithJSON(w, http.StatusOK, "Item added to cart successfully")
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

func (server *Server) UpdateCart(w http.ResponseWriter, r *http.Request) {
	cartID := GetShoppingCartID(w, r)
	cart, _ := GetShoppingCart(server.DB, cartID)

	for _, item := range cart.CartItems {
		qty, _ := strconv.Atoi(r.FormValue(item.ID))

		_, err := cart.UpdateItemQty(server.DB, item.ID, qty)
		if err != nil {
			jsonResponse := map[string]string{"error": "Error updating the cart"}
			sendJSONResponse(w, jsonResponse, http.StatusInternalServerError)
			return
		}
	}

	jsonResponse := map[string]string{"message": "Cart updated successfully"}
	sendJSONResponse(w, jsonResponse, http.StatusOK)
}

func (server *Server) RemoveItemByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["id"] == "" {
		jsonResponse := map[string]string{"error": "Invalid item ID"}
		sendJSONResponse(w, jsonResponse, http.StatusBadRequest)
		return
	}

	cartID := GetShoppingCartID(w, r)
	cart, _ := GetShoppingCart(server.DB, cartID)

	err := cart.RemoveItemByID(server.DB, vars["id"])
	if err != nil {
		jsonResponse := map[string]string{"error": "Error removing the item from the cart"}
		sendJSONResponse(w, jsonResponse, http.StatusInternalServerError)
		return
	}

	jsonResponse := map[string]string{"message": "Item removed from the cart"}
	sendJSONResponse(w, jsonResponse, http.StatusOK)
}
