package product

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
	
	"github.com/gorilla/mux"

	"leandroudala/foodstore/db"
	service "leandroudala/foodstore/services/product"
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

type errorParam interface {
	Error() string
}

func throwError(w http.ResponseWriter, err errorParam) {
	r := db.Status{400, err.Error()}
	w.WriteHeader(400)
	log.Println(r.Message)
	json.NewEncoder(w).Encode(r)
}

// GetProducts returns a list o products
func GetProducts(w http.ResponseWriter, r *http.Request) {
	list, err := service.GetProducts()

	if err != nil {
		throwError(w, err)
	} else {
		json.NewEncoder(w).Encode(list)
	}
}

// GetProduct returns a single item by ID
func GetProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	product, err := service.GetProduct(&id)

	if err != nil {
		throwError(w, err)
	}
	if product.ID == 0 {
		w.WriteHeader(http.StatusNotFound)	// 404 Not Found
	}

	json.NewEncoder(w).Encode(product)
}

// Create a new record to product's table
func Create(w http.ResponseWriter, r *http.Request) {
	product, err := service.Create(r.Body)

	if err != nil {
		throwError(w, err)
	}

	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(product)
}

// Update a product
func Update(w http.ResponseWriter, r *http.Request) {
	// preparing data
	product, err := service.Update(r.Body)
	if err != nil {
		throwError(w, err)
	}

	json.NewEncoder(w).Encode(product)
}

// Delete a product
func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	conn := db.GetConn()
	tx, err := conn.Begin()
	if err != nil {
		throwError(w, err)
	}

	stmt, err := tx.Prepare("delete from product where id = ?")
	if err != nil {
		throwError(w, err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)

	if err != nil {
		throwError(w, err)
	}

	rowsAffected, _ := result.RowsAffected()
	tx.Commit()
	if rowsAffected > 0 {
		w.WriteHeader(http.StatusOK) // 200
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
	}
}

func generateHash(b []byte) string {
	hash := sha1.New()
	hash.Write(b)
	return hex.EncodeToString(hash.Sum(nil))
}

// Upload receives a file
func Upload(w http.ResponseWriter, r *http.Request) {
	// reading uploaded file
	// var product Product
	
	file, handler, err := r.FormFile("image")
	if err != nil {
		throwError(w, err)
	}
	defer file.Close()

	// check if file is bigger than 5 MB
	if handler.Size > (5 << 20) {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	
	r.ParseForm()
	for i, j :=range r.Form {
		log.Println("\n", i,  " = ", j)
	}
	// log.Println(r.Form["description"][0])

	log.Println("name:", name)
	log.Println("description:", description)

	// generate extension
	ext := filepath.Ext(handler.Filename)
	// generating filename
	dt := time.Now().Format("20060102-150405.000")
	filename := "./uploads/" + generateHash([]byte(dt+handler.Filename)) + ext
	log.Println(filename)
	// creating output file
	// output, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755)
	// if err != nil {
	// 	throwError(w, err)
	// }
	// defer output.Close()

	// if _, err := io.Copy(output, file); err != nil {
	// 	throwError(w, err)
	// }
	w.WriteHeader(http.StatusOK)
}
