package product

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/renegmed/learn-go-webservice/inventoryservice/database"

	_ "github.com/go-sql-driver/mysql"
)

type httpResponse struct {
	Body        string
	StatusCode  int
	ContentType string
}

func removeAllProducts() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM products`)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func TestHandler(t *testing.T) {

	database.SetupDatabase()
	err := removeAllProducts()
	if err != nil {
		assertNoError(t, err, "error from database deleting all products,")
	}

	t.Run("it should be able to request for empty product list", func(t *testing.T) {

		wantJSON := `[]`

		request, err := newHandlerGetRequest("/products")
		assertNoError(t, err, "error while creating request,")

		response := httptest.NewRecorder()

		productsHandler(response, request)

		httpResp, err := getHttpResponse(response)

		assertNoError(t, err, "error from handler response,")
		assertStatusCode(t, httpResp.StatusCode, 200)
		assertContentType(t, httpResp.ContentType, "application/json")
		assertResponseJsonBody(t, httpResp.Body, wantJSON)
	})

	t.Run("it should able to post a product in json format", func(t *testing.T) {

		postJSON := `
{  
	"manufacturer":"Small Box Company",
	"sku":"4hs1j90JKL",
	"upc":"42465000000",
	"pricePerUnit":"9.99",
	"quantityOnHand":18,
	"productName":"Sprocket"
}`
		request, err := newHandlerPostRequestWithJson("/products", postJSON)
		assertNoError(t, err, "error while posting a request,")

		response := httptest.NewRecorder()

		productsHandler(response, request)

		httpResp, err := getHttpResponse(response)
		assertNoError(t, err, "error from handler response,")
		assertStatusCode(t, httpResp.StatusCode, 201)
	})

	t.Run("to verify added product from the product list", func(t *testing.T) {
		request, err := newHandlerGetRequest("/products")
		assertNoError(t, err, "error while creating request,")

		response := httptest.NewRecorder()

		productsHandler(response, request)

		httpResp, err := getHttpResponse(response)
		assertNoError(t, err, "error from handler response, ")

		assertStatusCode(t, httpResp.StatusCode, 200)
		assertContentType(t, httpResp.ContentType, "application/json")

		var products []Product

		err = json.Unmarshal([]byte(httpResp.Body), &products)
		if err != nil {

			assertNoError(t, fmt.Errorf("Error on unmarshalling products"), "")
			t.Fail()
		}
		if len(products) != 1 {
			assertNoError(t, fmt.Errorf("Product list should have only 1 product not %d.", len(products)), "")
			t.Fail()
		}
		product := products[0]
		if product.Manufacturer != "Small Box Company" {
			assertNoError(t, fmt.Errorf("Product manufacturer should be '%s' not '%s'", "Small Box Company", product.Manufacturer), "")
			t.Fail()
		}

		if product.PricePerUnit != fmt.Sprintf("%.2f", 9.99) {
			assertNoError(t, fmt.Errorf("Product price per unit %s not %s", fmt.Sprintf("%.2f", 9.99), product.PricePerUnit), "")
			t.Fail()
		}
	})

	t.Run("it should be able to request a particular product using product id", func(t *testing.T) {
		product := getFirstProduct(t)

		log.Println("product.ProductID:", *product.ProductID)

		request, err := newHandlerGetRequest("/products/" + fmt.Sprintf("%d", *product.ProductID))
		assertNoError(t, err, "error while requestin for a product,")

		response := httptest.NewRecorder()

		productHandler(response, request)

		httpResp, err := getHttpResponse(response)
		if err != nil {
			t.Fatalf("error from handler response, %v", err)
		}

		assertStatusCode(t, httpResp.StatusCode, 200)
		assertContentType(t, httpResp.ContentType, "application/json")

		var prod Product
		err = json.Unmarshal([]byte(httpResp.Body), &prod)
		if err != nil {
			assertNoError(t, fmt.Errorf("Error on unmarshalling product"), "")
			t.Fail()
		}

		if *product.ProductID != *prod.ProductID {
			assertNoError(t, fmt.Errorf("Product ID should be %s not %s", fmt.Sprintf("%d", product.ProductID), fmt.Sprintf("%d", prod.ProductID)), "")
			t.Fail()
		}

	})

	t.Run("it should be able to DELETE a product", func(t *testing.T) {
		product := getFirstProduct(t)

		request, err := newHandlerDeleteRequest("/products/" + fmt.Sprintf("%d", *product.ProductID))
		assertNoError(t, err, "error while requestin for a product,")

		response := httptest.NewRecorder()

		productHandler(response, request)

		httpResp, err := getHttpResponse(response)
		assertNoError(t, err, "error from handler response,")

		assertStatusCode(t, httpResp.StatusCode, 200)

		// Verify

		afterDeleteProductListJson := productList(t)

		assertResponseJsonBody(t, afterDeleteProductListJson, "[]")
	})
}

func getFirstProduct(t *testing.T) Product {

	request, err := newHandlerGetRequest("/products")
	assertNoError(t, err, "error while creating request,")

	response := httptest.NewRecorder()

	productsHandler(response, request)

	httpResp, err := getHttpResponse(response)
	assertNoError(t, err, "error from handler response, ")

	assertStatusCode(t, httpResp.StatusCode, 200)
	assertContentType(t, httpResp.ContentType, "application/json")

	var products []Product

	err = json.Unmarshal([]byte(httpResp.Body), &products)
	if err != nil {
		//log.Println("Error on unmarshalling products")
		assertNoError(t, fmt.Errorf("Error on unmarshalling products"), "")
		t.Fail()
	}
	if len(products) < 1 {
		assertNoError(t, fmt.Errorf("Product list should have one or more products not %d.", len(products)), "")
		t.Fail()
	}
	return products[0]
}

func productList(t *testing.T) string {
	request, err := newHandlerGetRequest("/products")
	if err != nil {
		t.Fatalf("error while creating request, %v", err)
	}
	response := httptest.NewRecorder()

	productsHandler(response, request)

	httpResp, err := getHttpResponse(response)
	if err != nil {
		t.Fatalf("error from handler response, %v", err)
	}

	assertStatusCode(t, httpResp.StatusCode, 200)

	return httpResp.Body
}

func getHttpResponse(response *httptest.ResponseRecorder) (httpResponse, error) {
	httpResponse := httpResponse{}
	resp := response.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return httpResponse, fmt.Errorf("error from response, %v", err)
	}

	httpResponse.Body = string(body)
	httpResponse.StatusCode = resp.StatusCode
	httpResponse.ContentType = resp.Header.Get("content-type")

	return httpResponse, nil
}

func newHandlerGetRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprint(url), nil)
	return req, err
}
func newHandlerPostRequestWithJson(url string, data string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprint(url), bytes.NewBuffer([]byte(data)))
	return req, err
}

func newHandlerPutRequestWithJson(url string, data string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprint(url), bytes.NewBuffer([]byte(data)))
	return req, err
}

func newHandlerDeleteRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprint(url), nil)
	return req, err
}

func assertNoError(t *testing.T, err error, message string) {
	if err != nil {
		t.Fatalf(message+" %v", err)
	}
}

func assertStatusCode(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("Response status code should be '%d' not '%d'", want, got)
	}
}

func assertContentType(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("Response content type should be '%s' not '%s'", want, got)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("Got response body  '%s' want '%s'", got, want)
	}
}

func assertResponseJsonBody(t *testing.T, got, want string) {

	got = cleanString(got)
	want = cleanString(want)
	if got != want {
		t.Errorf("Got response body \n'%s' \nwant \n'%s'", got, want)
	}
}

func cleanString(s string) string {
	ret := strings.ReplaceAll(s, "\t", "")
	ret = strings.ReplaceAll(ret, "\n", "")
	return ret
}
