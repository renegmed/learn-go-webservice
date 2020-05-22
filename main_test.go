package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type httpResponse struct {
	Body       string
	StatusCode int
}

func TestServerHTTP(t *testing.T) {
	fooHandler := fooHandler{"Hello World!!!"}

	request, err := newHandlerRequest("/foo")
	if err != nil {
		t.Fatalf("error while creating request, %v", err)
	}
	response := httptest.NewRecorder()

	fooHandler.ServeHTTP(response, request)

	httpResp, err := getHttpResponse(response)
	if err != nil {
		t.Fatalf("error from handler response, %v", err)
	}

	assertStatusCode(t, httpResp.StatusCode, 200)

	assertResponseBody(t, httpResp.Body, fooHandler.Message)
	if httpResp.Body != fooHandler.Message {
		t.Errorf("Got response body  '%s' want '%s'", httpResp.Body, fooHandler.Message)
	}
}

func TestGETHandler(t *testing.T) {
	want := "bar called"

	request, err := newHandlerRequest("/bar")
	if err != nil {
		t.Fatalf("error while creating request, %v", err)
	}
	response := httptest.NewRecorder()

	barHandler(response, request)

	httpResp, err := getHttpResponse(response)
	if err != nil {
		t.Fatalf("error from handler response, %v", err)
	}

	assertStatusCode(t, httpResp.StatusCode, 200)

	if httpResp.Body != want {
		t.Errorf("Got response body  '%s' want '%s'", httpResp.Body, want)
	}

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

	return httpResponse, nil
}
func newHandlerRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprint(url), nil)
	return req, err
}

func assertStatusCode(t *testing.T, statusCode, want int) {
	if statusCode != want {
		t.Errorf("Response status code should be '%d' not '%d'", want, statusCode)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("Got response body  '%s' want '%s'", got, want)
	}
}
