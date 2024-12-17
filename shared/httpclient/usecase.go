package httpclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func NewRequest(url, method string, payload interface{}, headers map[string]string) (request *http.Request, err error) {
	// decode the payload
	marshalled, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// create request object
	request, err = http.NewRequest(method, url, bytes.NewBuffer(marshalled))
	if err != nil {
		return
	}

	// apply headers
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	return
}

func DoRequest(url, method string, reqPayload interface{}, headers map[string]string, respHolderPTR interface{}) (statusCode int, err error) {
	statusCode = http.StatusInternalServerError

	request, err := NewRequest(url, method, reqPayload, headers)
	if err != nil {
		return
	}

	// do the http request
	response, err := Request(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	// get the status code
	statusCode = response.StatusCode

	// parse resp body
	err = ParseResponseBodyWithReference(response, &respHolderPTR)
	return
}

// Request function to do http request, default 30 second timeout
func Request(request *http.Request) (*http.Response, error) {
	netClient.Timeout = clientTimeout * time.Second
	resp, err := netClient.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RequestWithTimeout to do http request with user defined timeout
func RequestWithTimeout(request *http.Request, timeoutInSecond int) (*http.Response, error) {
	netClient.Timeout = time.Duration(timeoutInSecond) * time.Second
	resp, err := netClient.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ParseResponseBody is a func to parse json response body to interface
func ParseResponseBody(response *http.Response) (data interface{}, err error) {
	// Get Response Body
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	// decode response data to get the accessToken
	err = json.Unmarshal(responseData, &data)
	return
}

// ParseResponseBodyWithReference is a func to parse json response body to a placeholder
func ParseResponseBodyWithReference(response *http.Response, responsePlaceholder interface{}) (err error) {
	// Get Response Body
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	// decode response data to get the accessToken
	err = json.Unmarshal(responseData, &responsePlaceholder)
	return
}

// GetRoutePattern returns endpoint pattern of request
func GetRoutePattern(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())

	if pattern := rctx.RoutePattern(); pattern != "" {
		// Pattern is already available
		return pattern
	}
	routePath := r.URL.Path

	if r.URL.RawPath != "" {
		routePath = r.URL.RawPath
	}

	tctx := chi.NewRouteContext()

	if !rctx.Routes.Match(tctx, r.Method, routePath) {
		// No matching pattern, so just return the request path.
		// Depending on your use case, it might make sense to
		// return an empty string or error here instead
		return routePath
	}

	// tctx has the updated pattern, since Match mutates it
	return tctx.RoutePattern()

}
