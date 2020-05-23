package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type httpResponse struct {
	Body        string
	StatusCode  int
	ContentType string
}

const wantJSON = `[
{
	"productId":1,
	"manufacturer":"Johns-Jenkins",
	"sku":"p5z34vdS",
	"upc":"939581000000",
	"pricePerUnit":"497.45",
	"quantityOnHand":9703,
	"productName":"sticky note"
},
{
	"productId":2,
	"manufacturer":"Hessel, Schimmel and Feeney",
	"sku":"i7v300kmx",
	"upc":"740979000000",
	"pricePerUnit":"282.29",
	"quantityOnHand":9217,
	"productName":"leg warmers"
},
{
	"productId":3,
	"manufacturer":"Swaniawski, Bartoletti and Bruen",
	"sku":"q0L657ys7",
	"upc":"11173000000",
	"pricePerUnit":"436.26",
	"quantityOnHand":5905,
	"productName":"lamp shade"
}
]`

const newWantJSON = `
{
	"productId":4,
	"manufacturer":"Small Box Company",
	"sku":"4hs1j90JKL",
	"upc":"42465000000",
	"pricePerUnit":"9.99",
	"quantityOnHand":18,
	"productName":"Sprocket"
}`

func TestHandler(t *testing.T) {
	t.Run("it should be able to request for product list in json format", func(t *testing.T) {
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
		if err != nil {
			t.Fatalf("error while posting a request, %v", err)
		}
		response := httptest.NewRecorder()

		productsHandler(response, request)

		httpResp, err := getHttpResponse(response)
		if err != nil {
			t.Fatalf("error from handler response, %v", err)
		}

		assertStatusCode(t, httpResp.StatusCode, 201)
	})

	t.Run("to verify added product from the product list", func(t *testing.T) {
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
		assertContentType(t, httpResp.ContentType, "application/json")
		wantData := strings.ReplaceAll(strings.ReplaceAll(wantJSON, "[", ""), "]", ",") + newWantJSON
		wantData = "[" + wantData + "]"
		assertResponseJsonBody(t, httpResp.Body, wantData)
	})

	t.Run("it should be able to request a particular product using product id", func(t *testing.T) {
		wantProductJSON := `
		{
			"productId":2,
			"manufacturer":"Hessel, Schimmel and Feeney",
			"sku":"i7v300kmx",
			"upc":"740979000000",
			"pricePerUnit":"282.29",
			"quantityOnHand":9217,
			"productName":"leg warmers"
		}`

		request, err := newHandlerGetRequest("/products/2")
		if err != nil {
			t.Fatalf("error while requestin for a product, %v", err)
		}
		response := httptest.NewRecorder()

		productHandler(response, request)

		httpResp, err := getHttpResponse(response)
		if err != nil {
			t.Fatalf("error from handler response, %v", err)
		}

		assertStatusCode(t, httpResp.StatusCode, 200)
		assertContentType(t, httpResp.ContentType, "application/json")

		assertResponseJsonBody(t, httpResp.Body, wantProductJSON)
	})

	// t.Run("it should be able to update an existing product using PUT method", func(t *testing.T) {

	// 	modifiedProductJSON := `
	// 	{
	// 	"productId":4,
	// 	"manufacturer":"Small Box Company",
	// 	"sku":"4hs1j90JKL",
	// 	"upc":"42465000000",
	// 	"pricePerUnit":"9.99",
	// 	"quantityOnHand":215,
	// 	"productName":"Sprocket"
	// 	}`

	// 	request, err := newHandlerPutRequestWithJson("/products/4", modifiedProductJSON)
	// 	if err != nil {
	// 		t.Fatalf("error on product update request, %v", err)
	// 	}
	// 	response := httptest.NewRecorder()

	// 	productsHandler(response, request)

	// 	httpResp, err := getHttpResponse(response)
	// 	if err != nil {
	// 		t.Fatalf("error from handler response, %v", err)
	// 	}

	// 	assertStatusCode(t, httpResp.StatusCode, 200)

	// 	// verify the product updates

	// 	request, err = newHandlerGetRequest("/products/4")
	// 	if err != nil {
	// 		t.Fatalf("error while requesting for a product, %v", err)
	// 	}
	// 	response = httptest.NewRecorder()

	// 	productHandler(response, request)

	// 	httpResp, err = getHttpResponse(response)
	// 	if err != nil {
	// 		t.Fatalf("error from handler response, %v", err)
	// 	}

	// 	assertStatusCode(t, httpResp.StatusCode, 200)
	// 	assertContentType(t, httpResp.ContentType, "application/json")
	// 	assertResponseJsonBody(t, httpResp.Body, modifiedProductJSON)
	// })

	t.Run("it should be able to DELETE a product", func(t *testing.T) {

		//beforeDeleteProductListJson := getProductList(t)
		//log.Println("^^^^^^^ before delete:\n", beforeDeleteProductListJson)

		request, err := newHandlerDeleteRequest("/products/4")
		if err != nil {
			t.Fatalf("error while requestin for a product, %v", err)
		}
		response := httptest.NewRecorder()

		productHandler(response, request)

		httpResp, err := getHttpResponse(response)
		if err != nil {
			t.Fatalf("error from handler response, %v", err)
		}

		assertStatusCode(t, httpResp.StatusCode, 200)

		// Verify

		afterDeleteProductListJson := getProductList(t)
		// log.Println("^^^^^^^ after delete:\n", afterDeleteProductListJson)
		assertResponseJsonBody(t, afterDeleteProductListJson, wantJSON)
	})

}

func getProductList(t *testing.T) string {
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
	assertResponseJsonBody(t, httpResp.Body, wantJSON)

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
