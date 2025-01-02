package maputil

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

// HTTPResponseMap is a map representation of the response
//
// Example:
//
//	{
//		"request": {
//			"raw": "GET /foo/bar?one=testone&two=testtwo HTTP/1.1\nSome-Header: test\n\nfoo=1&bar=test"
//		},
//		"response": {
//			"content_length": "2034",
//			"status_code": "200",
//			"body": "<html>....</html>"
//			"raw": "HTTP/2 200 OK\nSome-Header: test\n\n<html>....</html>"
//		},
//		"headers": {
//			"Some-Header": "test"
//		}
//	}
type HTTPResponseMap map[string]interface{}

// Response returns http.Response instance by loading from raw request
func (hrm HTTPResponseMap) Response() (*http.Response, error) {
	reqObj := hrm["request"].(map[string]interface{})

	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(reqObj["raw_without_body"].(string))))
	if err != nil {
		return nil, err
	}

	respObj := hrm["response"].(map[string]interface{})
	return http.ReadResponse(bufio.NewReader(strings.NewReader(respObj["raw"].(string))), req)
}

// ContentLength returns response's content length
func (hrm HTTPResponseMap) ContentLength() int {
	respObj := hrm["response"].(map[string]interface{})
	return respObj["content_length"].(int)
}

// StatusCode returns response's status code
func (hrm HTTPResponseMap) StatusCode() int {
	respObj := hrm["response"].(map[string]interface{})
	return respObj["status_code"].(int)
}

// Body returns response's body as string
func (hrm HTTPResponseMap) Body() string {
	respObj := hrm["response"].(map[string]interface{})
	return respObj["body"].(string)
}

// Headers returns all header as http.Header instance
func (hrm HTTPResponseMap) Headers() http.Header {
	respObj := hrm["response"].(map[string]interface{})
	return respObj["headers"].(http.Header)
}

func NewHTTPResponseMap(resp *http.Response) (HTTPResponseMap, error) {
	mm := make(map[string]interface{})

	response := make(map[string]interface{})

	// Setup root variables
	response["proto"] = resp.Proto
	response["status_code"] = resp.StatusCode
	response["status_text"] = http.StatusText(resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	response["body"] = string(body)
	response["content_length"] = len(body)

	r, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	response["raw"] = string(r)

	// Setup headers
	response["headers"] = resp.Header

	// Setup request
	mReq, err := NewHTTPRequestMap(resp.Request)
	if err != nil {
		return nil, err
	}
	requestObject := mReq["request"].(map[string]interface{})
	delete(requestObject, "body")
	reqRaw := requestObject["raw"].(string)
	delete(requestObject, "raw")
	requestObject["raw_without_body"] = reqRaw

	mm["response"] = response
	mm["request"] = requestObject

	return mm, nil
}
