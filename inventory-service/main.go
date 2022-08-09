package main

import (
	"log"
	"net/http"

	"github.com/jordbick/Golang/inventory-service/database"
	"github.com/jordbick/Golang/inventory-service/product"
	"github.com/jordbick/Golang/inventory-service/receipt"

	// use underscore _ because we're not going to referencing the driver explicitly, just importing it for its side effects
	// and tin this case because we need the driver in order for the Go SQL package to work with our database
	_ "github.com/go-sql-driver/mysql"
)

// Need to add a new apiBasePath variable to pass into the SetupRoutes function
// /api for our base path
const basePath = "/api"

func main() {
	// call the function to create our DB variable
	database.SetupDatabase()
	product.SetupRoutes(basePath)
	receipt.SetupRoutes(basePath)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
