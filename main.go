package main

import (
	"log"
	"net/http"
	"flag"

	"github.com/gorilla/mux"

	"leandroudala/foodstore/controller/product"
	"leandroudala/foodstore/db"
	"leandroudala/foodstore/router"
)
// port that software is listening
var port *string

func init() {
	log.Println("Init server")
	sqlCommand := flag.String("sql", "", "Inform a valid SQL Command inside quotes.")
	port = flag.String("port", ":8000", "Inform a valid port number.")
	flag.Parse()

	for i, j := range flag.Args() {
		log.Println(i, j)
	}

	// check if it received a valid SQL Command
	if *sqlCommand != "" {
		db.ExecQuery(*sqlCommand)
	}
}

func main() {
	log.Println("main()")

	// connecting to DB
	if db.CheckDB() == false {
		db.CreateTables(
			product.SQL,
		)
	}

	// creating routes
	mx := mux.NewRouter()
	router.GetRoutes(mx)
	log.Println("Listening on port", *port)
	log.Fatal(http.ListenAndServe(*port, mx))
}
