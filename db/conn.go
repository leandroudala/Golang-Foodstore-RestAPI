package db

import (
	"database/sql"
	"log"
	"os"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

// CheckDB prepare DB File
const filename = "./db/files/database.db"

// Status is a object to store errors
type Status struct {
	Status  int    `json:"id"`
	Message string `json:"string"`
}

// CheckDB checks if DB file exists
func CheckDB() bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	} else if os.IsNotExist(err) {
		log.Printf(`File "%s" doesnt exist. I will create that.\n`, filename)
	} else {
		log.Fatal(err)
	}
	return false
}

// GetConn returns database connection
func GetConn() (db *sql.DB) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err.Error())
	}
	return
}

// CreateTables is used for create tables
func CreateTables(script ...string) {
	db, err := sql.Open("sqlite3", filename)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, j := range script {
		log.Println("Creating table:", j)
		_, err = db.Exec(j)
		if err != nil {
			log.Printf("%q: %s\n", err, j)
		}
	}
}

func checkTable(tableName string) bool {
	// open DB Connection
	db, err := sql.Open("sqlite3", filename)
	log.Println("Creating file", filename)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// check if table exists
	sqlQuery := "SELECT name FROM sqlite_master WHERE type='table' AND name='" + tableName + "';"
	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Fatal(err)
		return false
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Table name:", name)
	}

	return true
}
