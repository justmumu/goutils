package maputil

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// HTTPRequestMap is a map representation of the request
//
// Example:
//
//	{
//		"request": {
//			"method": "POST",
//			"path": "/foo/bar",
//			"body": "foo=1&bar=test",
//			"raw": "GET /foo/bar?one=testone&two=testtwo HTTP/1.1\nSome-Header: test\n\nfoo=1&bar=test"
//		},
//		"headers": {
//			"Some-Header": "test"
//		},
//		"query_params": {
//			"one": "testone",
//			"two": "testtwo"
//		}
//	}
type HTTPRequestMap map[string]interface{}

// Request returns http.Request instance by loading from raw request
func (hrm HTTPRequestMap) Request() (*http.Request, error) {
	reqObj := hrm["request"].(map[string]interface{})
	return http.ReadRequest(bufio.NewReader(strings.NewReader(reqObj["raw"].(string))))
}

// Method returns request's method
func (hrm HTTPRequestMap) Method() string {
	reqObj := hrm["request"].(map[string]interface{})
	return reqObj["method"].(string)
}

// Path returns request's req.URL.Path
func (hrm HTTPRequestMap) Path() string {
	reqObj := hrm["request"].(map[string]interface{})
	return reqObj["path"].(string)
}

// Body returns request's body as string
func (hrm HTTPRequestMap) Body() string {
	reqObj := hrm["request"].(map[string]interface{})
	return reqObj["body"].(string)
}

// Headers returns all header as http.Header instance
func (hrm HTTPRequestMap) Headers() http.Header {
	reqObj := hrm["request"].(map[string]interface{})
	return reqObj["headers"].(http.Header)
}

// QueryParams returns all query params as url.Values instance
func (hrm HTTPRequestMap) QueryParams() url.Values {
	reqObj := hrm["request"].(map[string]interface{})
	urlObj := reqObj["url"].(map[string]interface{})
	return urlObj["query_params"].(url.Values)
}

func NewHTTPRequestMap(req *http.Request) (HTTPRequestMap, error) {
	mm := make(map[string]interface{})

	// Setup Root Variables
	request := make(map[string]interface{})
	request["method"] = req.Method
	request["proto"] = req.Proto

	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(body))
		request["body"] = string(body)
	} else {
		request["body"] = ""
	}

	reqdump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	request["raw"] = string(reqdump)

	// Setup URL data
	url := make(map[string]interface{})
	url["scheme"] = req.URL.Scheme
	url["user"] = req.URL.User.String()
	url["hostname"] = req.URL.Hostname()
	url["port"] = req.URL.Port()
	url["path"] = req.URL.Path
	url["query_params"] = req.URL.Query()
	url["uri"] = req.URL.RequestURI()
	url["raw"] = req.URL.String()
	request["url"] = url

	// Setup headers
	request["headers"] = req.Header

	// Inject request to root
	mm["request"] = request

	return mm, nil
}
