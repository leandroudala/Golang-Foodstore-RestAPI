package services

import (
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
func GetProduct(id int) (model.Product, error) {
	var product model.Product

	conn := db.GetConn()
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		return product, err
	}

	stmt, err := tx.Prepare("select id, name, description, price from product where id = ?")
	if err != nil {
		return product, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return product, err
	}
	defer rows.Close()

	
	if rows.Next() {
		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
		if err != nil {
			return product, err
		}		
	}

	return product, nil
}
