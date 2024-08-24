package openai_test

import (
	"github.com/skyscrapr/openai-sdk-go/openai"
	"net/url"
)

const test_api_key = "this-is-my-secure-apikey-do-not-steal!!"
const test_admin_key = "this-is-my-secure-adminkey-do-not-steal!!"
const test_organization_id = "this-is-my-organization-id"

func GetTestAuthToken() string {
	return test_api_key
}

func GetTestAdminToken() string {
	return test_admin_key
}

func NewTestClient(ts *TestServer) *openai.Client {
	client := openai.NewClient(test_api_key, test_admin_key)
	client.OrganizationID = test_organization_id
	if ts != nil {
		client.BaseURL, _ = url.Parse(ts.HTTPServer.URL)
	}
	return client
}
