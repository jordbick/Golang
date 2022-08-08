package database

import (
	"database/sql"
	"log"
)

// In the SQL package within go, there's a function called Open
// It takes a driver name for the database we're going to be using as well as the data source name (connection string for the DB)
// Returns a DB object and a possible error
// func Open(driverName, dataSourceName string) (*DB, error)

// The DB type is a handle to the database which is responsible for managing the connections within a connection pool
// The DB is going to be automatically opening and closing new connections as needed, will perform this in a thread-safe way

// variable to hold our DB connection object
// Exported so is capital letter
var DbConn *sql.DB

// In order to use the SQL package we also need a driver for the specific database we're going to use
// The driver itself isn't part of the Go starndard library
func SetupDatabase() {
	var err error
	// mysql is driver name
	// connection string = username:password@host:port/nameofDB
	DbConn, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/inventorydb")
	if err != nil {
		log.Fatal(err)
	}
}

// Interacting with DB
// Query function allows us to return one or more rows
// Takes in a query string and a variadic list of arguments that's intended to replace the parameters within the query
// Returns a Rows object (result of a SQL Query) and a possible error
// In order to consume the data in the Rows type we will need to iterate over it and use the Next method on the Rows
// Need to make sure we close our Rows object once we're done with it
// func (db *DB) Query(query string, args ...interface{}) (*Rows, error)

// In order to parse out the data from the Rows we use the method Scan
// Scan convertes the columns read from the database into our Go types
// func (rs *Rows) Scan(dest ...interface{}) error

// 				// db query with simple select statement
// results, err := db.Query(`select productId, manufacturer, sku from products`)
// if err != nil {
// 	log.Fatal(err)
// }
// defer results.Close()
// products := make([]Product, 0)
// // Iterate through the results and calling results.Next to advance our cursos for each iteration of the loop
// for results.Next(){
// 	var product Product
// 	// results.Scan method and pass in the specific fields in out new product variable that we want to set
// 	results.Scan(&product.ProductID, &product.Manufacturer, &product.Sku ...)
// 	products = append(products, product)
// }

// Another query method called QueryRow
// This method is only ever expected to return a single row
// If it returns more than one then the Rows.Scan method will grab the first row and discard the rest
// func (db *DB) QueryRow(query string, args ...interface{}) *Row

// The row type only has one mehtod attached to it which is Scan
// func (rs *Row) Scan(dest ...interface{}) error

// To download mysql driver for go
// go get -u github.com/go-sql-driver/mysql
// also add import to main.go file
// _ "github.com/go-sql-driver/mysql"
