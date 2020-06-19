package main

import (
	"log"
	"net/http"

	"github.com/renegmed/inventoryservice/receipt"

	"github.com/renegmed/inventoryservice/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/renegmed/inventoryservice/product"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	receipt.SetupRoutes(apiBasePath)
	product.SetupRoutes(apiBasePath)
	log.Println("Server started on port 5000...")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
