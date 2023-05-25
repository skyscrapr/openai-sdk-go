package openai

import (
	"testing"
)

func TestNewEndpoint(t *testing.T) {
	testEndpointPath := "testEndpointPath"
	testClient := NewClient("testapikey")
	e := newEndpoint(testClient, testEndpointPath)
	if e.BaseURL.String() != testClient.BaseURL.String() {
		t.Errorf("VendorsEndpoint BaseURL mismatch. Got %s. Want %s", e.BaseURL.String(), testClient.BaseURL.String())
	}
	if e.EndpointPath != testEndpointPath {
		t.Errorf("VendorsEndpoint EndpointPath mismatch. Got %s. Expected %s", e.EndpointPath, testEndpointPath)
	}
}
