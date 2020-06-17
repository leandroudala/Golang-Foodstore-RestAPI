package product

import (
	"encoding/json"
	"leandroudala/foodstore/db"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	// SQL is SQL Query for create a new table
	SQL = `
		CREATE TABLE product (
			id integer primary key autoincrement,
			name varchar(75) not null,
			description text null,
			price float null
		)
	`
)

// Product stores data about product
type Product struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

type errorParam interface {
	Error() string
}

func throwError(w http.ResponseWriter, err errorParam) {
	r := db.Status{400, err.Error()}
	log.Println(r.Message)
	json.NewEncoder(w).Encode(r)
}

// GetProducts return a list o products
func GetProducts(w http.ResponseWriter, r *http.Request) {
	conn := db.GetConn()

	defer conn.Close()
	rows, err := conn.Query("select id, name, description, price from product")
	if err != nil {
		throwError(w, err)
	}

	list := make([]Product, 0)
	defer rows.Close()
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price)
		if err != nil {
			throwError(w, err)
		}
		list = append(list, p)
	}

	json.NewEncoder(w).Encode(list)
}

// GetProduct returns a single item by ID
func GetProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	conn := db.GetConn()
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		throwError(w, err)
	}

	stmt, err := tx.Prepare("select id, name, description, price from product where id = ?")
	if err != nil {
		throwError(w, err)
	}
	rows, err := stmt.Query(id)
	if err != nil {
		throwError(w, err)
	}
	defer rows.Close()

	var product Product
	if rows.Next() {
		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
		if err != nil {
			throwError(w, err)
		}

		json.NewEncoder(w).Encode(product)
	} else {
		w.WriteHeader(404)
	}
}

// Create inserts a new record to product's table
func Create(w http.ResponseWriter, r *http.Request) {
	// preparing data
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Println(err)
	}

	conn := db.GetConn()
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		log.Println("Error when inserting:", err.Error)
		return
	}
	stmt, err := tx.Prepare("insert into product (name, description, price) values (?, ?, ?)")
	if err != nil {
		throwError(w, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(product.Name, product.Description, product.Price)

	if err != nil {
		throwError(w, err)
	}
	tx.Commit()

	id, err := res.LastInsertId()
	if err != nil {
		throwError(w, err)
	}
	product.ID = uint(id)
	json.NewEncoder(w).Encode(product)
}
