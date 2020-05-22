package main

import (
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

func TestHandler(t *testing.T) {
	wantJSON := `[
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
	"productId":12,
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
	t.Run("it should able to request for product list in json format", func(t *testing.T) {
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

func assertStatusCode(t *testing.T, statusCode, want int) {
	if statusCode != want {
		t.Errorf("Response status code should be '%d' not '%d'", want, statusCode)
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

	// if !reflect.DeepEqual(got, want) { //got != want {
	// 	t.Errorf("Got response body  '%s' want '%s'", got, want)
	// }
	got = cleanString(got)
	want = cleanString(want)
	if got != want {
		t.Errorf("Got response body  \n'%s' want \n'%s'", got, want)
	}
}

func cleanString(s string) string {
	ret := strings.ReplaceAll(s, "\t", "")
	ret = strings.ReplaceAll(ret, "\n", "")
	return ret
}
