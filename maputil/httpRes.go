package maputil

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
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
type HTTPResponseMap map[string]map[string]string

// Response returns http.Response instance by loading from raw request
func (hrm HTTPResponseMap) Response() (*http.Response, error) {
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(bytes.NewBufferString(hrm["request"]["raw"]).Bytes())))
	if err != nil {
		return nil, err
	}
	return http.ReadResponse(bufio.NewReader(bytes.NewReader(bytes.NewBufferString(hrm["response"]["raw"]).Bytes())), req)
}

// ContentLength returns response's content length
func (hrm HTTPResponseMap) ContentLength() int64 {
	val, _ := strconv.ParseInt(hrm["response"]["content_length"], 10, 16)
	return val
}

// StatusCode returns response's status code
func (hrm HTTPResponseMap) StatusCode() int {
	val, _ := strconv.ParseInt(hrm["response"]["status_code"], 10, 16)
	return int(val)
}

// Body returns response's body as string
func (hrm HTTPResponseMap) Body() string {
	return hrm["response"]["body"]
}

// Headers returns all header as http.Header instance
func (hrm HTTPResponseMap) Headers() http.Header {
	returnVal := make(map[string][]string)

	for k, vv := range hrm["headers"] {
		returnVal[k] = strings.Split(vv, multiValueSeparator)
	}
	return returnVal
}

func NewHTTPResponseMap(resp *http.Response) (HTTPResponseMap, error) {
	m := make(map[string]map[string]string)

	request := make(map[string]string)
	req, err := httputil.DumpRequest(resp.Request, true)
	if err != nil {
		return nil, err
	}
	request["raw"] = string(req)
	m["request"] = request

	// Setup response
	response := make(map[string]string)
	response["content_length"] = fmt.Sprintf("%d", resp.ContentLength)
	response["status_code"] = fmt.Sprintf("%d", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	response["body"] = string(body)

	r, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	response["response"] = string(r)
	m["response"] = response

	hm := make(map[string]string)
	for k, v := range resp.Header {
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

	return m, nil
}
