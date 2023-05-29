package openai_test

import (
	"net/url"
	"github.com/skyscrapr/openai-sdk-go/openai"
)

const test_auth_token = "this-is-my-secure-token-do-not-steal!!"

func GetTestAuthToken() string {
	return test_auth_token
}


func NewTestClient(ts *TestServer) *openai.Client {
	client := openai.NewClient(test_auth_token)
	client.BaseURL, _ = url.Parse(ts.HTTPServer.URL)
	return client
}
