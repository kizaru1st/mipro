package controllers

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/shopspring/decimal"

	_ "github.com/shopspring/decimal"

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
		existCart, _ = cart.CreateCart(db, cartID)
	}

	_, _ = existCart.CalculateCart(db, cartID)

	updatedCart, _ := cart.GetCart(db, cartID)

	totalWeight := 0
	productModel := models.Product{}
	for _, cartItem := range updatedCart.CartItems {
		product, _ := productModel.FindByID(db, cartItem.ProductID)

		productWeight, _ := product.Weight.Float64()
		ceilWeight := math.Ceil(productWeight)

		itemWeight := cartItem.Qty * int(ceilWeight)

		totalWeight += itemWeight
	}

	updatedCart.TotalWeight = totalWeight

	return updatedCart, nil
}

func (server *Server) GetCart(w http.ResponseWriter, r *http.Request) {
	cartID := GetShoppingCartID(w, r)
	cart, _ := GetShoppingCart(server.DB, cartID)
	items, _ := cart.GetItems(server.DB, cartID)

	provinces, err := server.GetProvinces()
	if err != nil {
		jsonResponse := map[string]string{"error": "Failed to retrieve provinces"}
		sendJSONResponse(w, jsonResponse, http.StatusInternalServerError)
		return
	}

	jsonResponse := map[string]interface{}{
		"cart":      cart,
		"items":     items,
		"provinces": provinces,
	}

	sendJSONResponse(w, jsonResponse, http.StatusOK)
}

func (server *Server) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	productID := r.FormValue("product_id")
	qty, _ := strconv.Atoi(r.FormValue("qty"))

	productModel := models.Product{}
	product, err := productModel.FindByID(server.DB, productID)
	if err != nil {
		jsonResponse := map[string]string{"error": "Product not found"}
		sendJSONResponse(w, jsonResponse, http.StatusNotFound)
		return
	}

	if qty > product.Stock {
		jsonResponse := map[string]string{"error": "Quantity exceeds available stock"}
		sendJSONResponse(w, jsonResponse, http.StatusBadRequest)
		return
	}

	var cart *models.Cart

	cartID := GetShoppingCartID(w, r)
	cart, _ = GetShoppingCart(server.DB, cartID)
	_, err = cart.AddItem(server.DB, models.CartItem{
		ProductID: productID,
		Qty:       qty,
	})
	if err != nil {
		jsonResponse := map[string]string{"error": "Failed to add item to cart"}
		sendJSONResponse(w, jsonResponse, http.StatusInternalServerError)
		return
	}

	jsonResponse := map[string]string{"message": "Item added to cart successfully"}
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
		jsonResponse := map[string]string{"message": "Product not found"}
		sendJSONResponse(w, jsonResponse, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(product); err != nil {
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
			jsonResponse := map[string]string{"error": "Failed to update the cart"}
			sendJSONResponse(w, jsonResponse, http.StatusInternalServerError)
			return
		}
	}

	jsonResponse := map[string]string{"message": "Cart updated successfully"}
	sendJSONResponse(w, jsonResponse, http.StatusOK)
}

func (server *Server) RemoveItem(w http.ResponseWriter, r *http.Request) {
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

func (server *Server) GetCitiesByProvince(w http.ResponseWriter, r *http.Request) {
	provinceID := r.URL.Query().Get("province_id")

	cities, err := server.GetCitiesByProvinceID(provinceID)
	if err != nil {
		jsonResponse := map[string]string{"error": "Failed to retrieve cities"}
		sendJSONResponse(w, jsonResponse, http.StatusInternalServerError)
		return
	}

	res := Result{Code: 200, Data: cities, Message: "Success"}
	result, err := json.Marshal(res)
	if err != nil {
		jsonResponse := map[string]string{"error": "Failed to marshal JSON"}
		sendJSONResponse(w, jsonResponse, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (server *Server) CalculateShipping(w http.ResponseWriter, r *http.Request) {
	origin := os.Getenv("API_ONGKIR_ORIGIN")
	destination := r.FormValue("city_id")
	courier := r.FormValue("courier")

	if destination == "" {
		jsonResponse := map[string]string{"error": "Invalid destination"}
		sendJSONResponse(w, jsonResponse, http.StatusBadRequest)
		return
	}

	cartID := GetShoppingCartID(w, r)
	cart, _ := GetShoppingCart(server.DB, cartID)

	shippingFeeOptions, err := server.CalculateShippingFee(models.ShippingFeeParams{
		Origin:      origin,
		Destination: destination,
		Weight:      cart.TotalWeight,
		Courier:     courier,
	})

	if err != nil {
		jsonResponse := map[string]string{"error": err.Error()}
		sendJSONResponse(w, jsonResponse, http.StatusInternalServerError)
		return
	}

	res := Result{Code: 200, Data: shippingFeeOptions, Message: "Success"}
	result, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (server *Server) ApplyShipping(w http.ResponseWriter, r *http.Request) {
	origin := os.Getenv("API_ONGKIR_ORIGIN")
	destination := r.FormValue("city_id")
	courier := r.FormValue("courier")
	shippingPackage := r.FormValue("shipping_package")

	cartID := GetShoppingCartID(w, r)
	cart, _ := GetShoppingCart(server.DB, cartID)

	if destination == "" {
		jsonResponse := map[string]string{"error": "Invalid destination"}
		sendJSONResponse(w, jsonResponse, http.StatusBadRequest)
		return
	}

	shippingFeeOptions, err := server.CalculateShippingFee(models.ShippingFeeParams{
		Origin:      origin,
		Destination: destination,
		Weight:      cart.TotalWeight,
		Courier:     courier,
	})

	if err != nil {
		jsonResponse := map[string]string{"error": "Invalid shipping calculation"}
		sendJSONResponse(w, jsonResponse, http.StatusInternalServerError)
		return
	}

	var selectedShipping models.ShippingFeeOption

	for _, shippingOption := range shippingFeeOptions {
		if shippingOption.Service == shippingPackage {
			selectedShipping = shippingOption
			continue
		}
	}

	type ApplyShippingResponse struct {
		TotalOrder  decimal.Decimal `json:"total_order"`
		ShippingFee decimal.Decimal `json:"shipping_fee"`
		GrandTotal  decimal.Decimal `json:"grand_total"`
		TotalWeight decimal.Decimal `json:"total_weight"`
	}

	var grandTotal decimal.Decimal

	grandTotal = cart.GrandTotal.Add(decimal.NewFromInt(int64(selectedShipping.Fee)))

	applyShippingResponse := ApplyShippingResponse{
		TotalOrder:  cart.GrandTotal,
		ShippingFee: decimal.NewFromInt(selectedShipping.Fee),
		GrandTotal:  grandTotal,
		TotalWeight: decimal.NewFromInt(int64(cart.TotalWeight)),
	}

	res := Result{Code: 200, Data: applyShippingResponse, Message: "Success"}
	result, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
