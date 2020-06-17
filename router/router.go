package router

import (
	"leandroudala/foodstore/controller/product"

	"github.com/gorilla/mux"
)

// GetRoutes is the main function that returns router
func GetRoutes(r *mux.Router) {
	r.HandleFunc("/product", product.Create).Methods("Post")
	r.HandleFunc("/product/{id}", product.GetProduct).Methods("Get")
	r.HandleFunc("/products", product.GetProducts).Methods("Get")
}
