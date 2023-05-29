package openai_test

import (
	"log"
	"net/http"
	"net/http/httptest"
)

type TestServer struct {
	HTTPServer *httptest.Server
	handlers map[string]handler
}
type handler func(w http.ResponseWriter, r *http.Request)

func NewTestServer() *TestServer {
	ts := TestServer{
		handlers: make(map[string]handler),
	}
	ts.HTTPServer = ts.enableHandlers()
	return &ts
}


func (ts *TestServer) RegisterHandler(path string, handler handler) {
	ts.handlers[path] = handler
}

// OpenAITestServer Creates a mocked OpenAI server which can pretend to handle requests during testing.
func (ts *TestServer) enableHandlers() *httptest.Server {
	return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("received request at path %q\n", r.URL.Path)

		// check auth
		if r.Header.Get("Authorization") != "Bearer " + GetTestAuthToken() {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handlerCall, ok := ts.handlers[r.URL.Path]
		if !ok {
			http.Error(w, "the resource path doesn't exist", http.StatusNotFound)
			return
		}
		handlerCall(w, r)
	}))
}
