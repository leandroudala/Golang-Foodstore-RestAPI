package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"leandroudala/foodstore/controller/product"
	"leandroudala/foodstore/db"
	"leandroudala/foodstore/router"
)

func main() {
	// connecting to DB
	if db.CheckDB() == false {
		db.CreateTables(
			product.SQL,
		)
	}

	// creating routes
	mx := mux.NewRouter()
	router.GetRoutes(mx)
	log.Println("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", mx))
}
