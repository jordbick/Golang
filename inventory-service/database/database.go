package database

import (
	"database/sql"
	"log"
	"time"
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
	DbConn.SetMaxOpenConns(4)
	DbConn.SetMaxIdleConns(4)
	DbConn.SetConnMaxLifetime(60 * time.Second)
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

// Other SQL statements such as INSERT, UPDATE or DELETE
// Use DB.Exec
// Get back as SQL Result object
// func (rs *DB) Exec(query string, args ...interface{}) (Result, error)

// Result interface has two functions
// type Result interface {
// LastInsertId() (int64, error)
// RowsAffected() (int64, error)
// }

// DB type is responsible for manafing the connections to the database
// Includes handling multiple concurrent requests at any time
// It does this by creating a pool of connections which can be reused when finsihed with

// Can use methods to configure th behaviour of the DB object and how it manages the connections to the DB

// Conection Max Lifetime - Sets the max amount of time a connection may be used
// Max Idle Connections - Sets the max number of connections in the idle connection pool
// Max Open Connections - Sets the max number of open connections to the DB
// Set these because for many DBs there is a max limit of connections that the DB is capable of handling
// By setting these we can control the flow by limiting the number of connections to your DB

// When max number of open connections exhausted may get connection timeouts
// Need to use a context
// Context = Mechanisms that allows you to set a deadline, cancel a signal, or set other request-scoped values across API boundaries and between processed
// To do this we need to create a context object

// context.WithTimeout method accepts a context and a Go time object as input. context.Background() creates a blank context
// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// Change Query method to QueryContext and pass in context object
// results, err := db.QueryContext(ctx, `select ... from ...`)
// Now if the query exceeds the time limit that we set in the timeout, the call will cancel and return

// Also QueryRowContext and ExecContext methods

// File Uploads
// Could send a file using JSON by first encoding the file into a base64 and then including that as a property in the JSON
// Uses the decode string method which accepts a base64 string as input and returns a slice of bytes
// func (enc *Encoding) DecodeString(s string) ([]byte, error)

// Can also send file by using the HTTP multipart/form-data content type which allows us to use a HTML form to submit the raw binary data to our web service
// This method is more efficient
// Returns a multipart file object and a pointer to a FileHeader which contains information about the file
// func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)

// The File type implements the io.Reader function which allows us to read the file content
// The FileHeader type can get the Filename, size and MIME type of the file
// In order to implement this in our handler we need to first call ParseMultiPartForm on the request, passing in a size which will limit the amount of data from the request that is stored in memory
// If the size of the request exceeds the limit then it will store data in temporary files on disk
// func uploadFileHandler(w http.ResponseWriter, r *http.Request){
// r.ParseMultiPartForm(5 << 20) // 5Mb

// Then grab the field from the form matching the key that we pass in - Expecting the form to have a key called uploadFileName in the request
// file, handler, err := r.FormFile("uploadFileName")

// Then use OS package OpenFile func to create a new file at the given filepath
// f, err := os.OpenFile("./filepath/" + handler.Finename, os.O_WRONLY|os.O_CREATE, 0666)

// Then use the io.Copy func to copy the byte data that we read in from the request to the new file
// io.Copy(f, file)

// To handle downloads
// Use os.Open function
// func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
// filename = "gopher.png"
// file, err := os.Open(fileName)
// }

// To tell the client that this is a file that should be downloaded we can set the header content disposition to attachment and give it a file name
// w.Header.Set("Content-Disposition", "attachment; filename="+fileName)
