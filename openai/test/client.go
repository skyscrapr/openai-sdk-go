package openai_test

import (
	"github.com/skyscrapr/openai-sdk-go/openai"
	"net/url"
)

const test_auth_token = "this-is-my-secure-token-do-not-steal!!"

func GetTestAuthToken() string {
	return test_auth_token
}

func NewTestClient(ts *TestServer) *openai.Client {
	client := openai.NewClient(test_auth_token)
	if (ts != nil) {
		client.BaseURL, _ = url.Parse(ts.HTTPServer.URL)
	}
	return client
}
