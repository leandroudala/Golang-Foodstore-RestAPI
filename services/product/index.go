package services

import (
	"io"
	"encoding/json"

	"leandroudala/foodstore/db"
	model "leandroudala/foodstore/models/product"
)

// GetProducts connects to DB and return a list
func GetProducts() ([]model.Product, error) {
	list := make([]model.Product, 0)
	
	conn := db.GetConn()
	defer conn.Close()
	rows, err := conn.Query("select id, name, description, price from product")
	if err != nil {
		return nil, err
	}

	
	defer rows.Close()
	for rows.Next() {
		var p model.Product
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price)
		
		if err != nil {
			return nil, err
		}
		
		list = append(list, p)
	}

	return list, nil
}

type errorParam struct {
	HTTPStatus int
	Error string
}

// GetProduct returns a specific product by id
func GetProduct(id *int) (*model.Product, error) {
	var product model.Product

	conn := db.GetConn()
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("select id, name, description, price from product where id = ?")
	if err != nil {
		return &product, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return &product, err
	}
	defer rows.Close()

	
	if rows.Next() {
		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
		if err != nil {
			return &product, err
		}		
	}

	return &product, nil
}


// Create a new record to product's table
func Create(body io.Reader) (*model.Product, error) {
	var product model.Product
	err := json.NewDecoder(body).Decode(&product)
	if err != nil {
		return nil, err
	}

	conn := db.GetConn()
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare("insert into product (name, description, price, image) values (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(product.Name, product.Description, product.Price, product.Image)

	if err != nil {
		return nil, err
	}
	tx.Commit()

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	product.ID = uint(id)
	
	return &product, nil
}


// Update a product
func Update(body io.Reader) (*model.Product, error) {
	// preparing data
	var product model.Product
	err := json.NewDecoder(body).Decode(&product)
	if err != nil {
		return nil, err
	}

	conn := db.GetConn()
	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("update product set name = ?, description = ?, price = ? where id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Description, product.Price, product.ID)
	tx.Commit()
	if err != nil {
		return nil, err
	} else {
		return &product, nil
	}
}