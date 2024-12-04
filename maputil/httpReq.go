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
type HTTPRequestMap map[string]map[string]string

// Request returns http.Request instance by loading from raw request
func (hrm HTTPRequestMap) Request() (*http.Request, error) {
	return http.ReadRequest(bufio.NewReader(bytes.NewReader(bytes.NewBufferString(hrm["request"]["raw"]).Bytes())))
}

// Method returns request's method
func (hrm HTTPRequestMap) Method() string {
	return hrm["request"]["method"]
}

// Path returns request's req.URL.Path
func (hrm HTTPRequestMap) Path() string {
	return hrm["request"]["path"]
}

// Body returns request's body as string
func (hrm HTTPRequestMap) Body() string {
	return hrm["request"]["body"]
}

// Headers returns all header as http.Header instance
func (hrm HTTPRequestMap) Headers() http.Header {
	returnVal := make(map[string][]string)

	for k, vv := range hrm["headers"] {
		returnVal[k] = strings.Split(vv, multiValueSeparator)
	}
	return returnVal
}

// QueryParams returns all query params as url.Values instance
func (hrm HTTPRequestMap) QueryParams() url.Values {
	returnVal := make(map[string][]string)

	for k, vv := range hrm["query_params"] {
		returnVal[k] = strings.Split(vv, multiValueSeparator)
	}

	return returnVal
}

func NewHTTPRequestMap(req *http.Request) (HTTPRequestMap, error) {
	m := make(map[string]map[string]string)

	// Setup request data
	request := make(map[string]string)
	request["method"] = req.Method
	request["path"] = req.URL.Path

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))
	request["body"] = string(body)

	reqdump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	request["raw"] = string(reqdump)

	m["request"] = request

	// Setup headers
	hm := make(map[string]string)
	for k, v := range req.Header {
		var vv string
		for i, val := range v {
			vv += strings.TrimSpace(val)
			if i+1 < len(v) {
				vv += multiValueSeparator
			}
		}
		hm[strings.TrimSpace(k)] = vv
	}
	m["headers"] = hm

	// Setup query params
	qvm := make(map[string]string)
	for k, v := range req.URL.Query() {
		var vv string
		for i, val := range v {
			vv += strings.TrimSpace(val)
			if i+1 < len(v) {
				vv += multiValueSeparator
			}
		}
		qvm[strings.TrimSpace(k)] = vv
	}
	m["query_params"] = qvm

	return m, nil
}
