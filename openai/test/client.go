package openai_test

import (
	"github.com/skyscrapr/openai-sdk-go/openai"
	"net/url"
)

const test_api_key = "this-is-my-secure-apikey-do-not-steal!!"

func GetTestAuthToken() string {
	return test_api_key
}

func NewTestClient(ts *TestServer) *openai.Client {
	client := openai.NewClient(test_api_key)
	if ts != nil {
		client.BaseURL, _ = url.Parse(ts.HTTPServer.URL)
	}
	return client
}
