package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func testV1API(t *testing.T, method string, path string, body interface{}) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			t.Errorf("Expected request method ‘%s’, got ‘%s’", method, req.Method)
		}
		if req.URL.Path != path {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", path, req.URL.Path)
		}
		if body != nil {
			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(body)
			if err != nil {
				t.Errorf("Error converting body into JSON: %s", err)
			}
			rw.Write(buf.Bytes())
		} else {
			// TODO: Fix this. I don't think it's right
			rw.Write([]byte(`{}`))
		}
	}))
	return server
}

func CheckStructEqual(t *testing.T, got interface{}, expected interface{}) {
	if !reflect.DeepEqual(got, expected) {
		sgot, _ := json.MarshalIndent(got, "", "\t")
		sexpected, _ := json.MarshalIndent(expected, "", "\t")
		t.Errorf("Expected '%s', got '%s'", sexpected, sgot)
	}
}
