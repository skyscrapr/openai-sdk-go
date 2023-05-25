package openai

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func newTestClient(t *testing.T, s *httptest.Server) *Client {
	client := NewClient("testapikey")
	u, err := url.Parse(s.URL)
	if err != nil {
		t.Fail()
	}
	client.BaseURL = u
	return client
}
