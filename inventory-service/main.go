package main

import (
	"log"
	"net/http"

	"github.com/renegmed/inventoryservice/product"
)

const apiBasePath = "/api"

func main() {
	product.SetupRoutes(apiBasePath)
	log.Println("Server started on port 5000...")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
