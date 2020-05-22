package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/renegmed/inventoryservice/product"
)

type fooHandler struct {
	Message string
}

var productList []product.Product

// func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte(f.Message))
// }

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/json")
		w.Write(productsJson)
	}
}

func init() {
	productsJSON := `[
{
	"productId": 1,
	"manufacturer": "Johns-Jenkins",
	"sku": "p5z34vdS",
	"upc": "939581000000",
	"pricePerUnit": "497.45",
	"quantityOnHand": 9703,
	"productName": "sticky note"
},
{
	"productId": 12,
	"manufacturer": "Hessel, Schimmel and Feeney",
	"sku": "i7v300kmx",
	"upc": "740979000000",
	"pricePerUnit": "282.29",
	"quantityOnHand": 9217,
	"productName": "leg warmers"
},
{
	"productId": 3,
	"manufacturer": "Swaniawski, Bartoletti and Bruen",
	"sku": "q0L657ys7",
	"upc": "11173000000",
	"pricePerUnit": "436.26",
	"quantityOnHand": 5905,
	"productName": "lamp shade"
}
]`
	err := json.Unmarshal([]byte(productsJSON), &productList)
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	http.HandleFunc("/products", productsHandler)
	http.ListenAndServe(":5000", nil)
}
