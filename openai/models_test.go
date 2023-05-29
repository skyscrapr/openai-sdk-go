package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skyscrapr/openai-sdk-go/openai"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
)

// TestListModels Tests the models endpoint of the API using the mocked server.
func TestListModels(t *testing.T) {
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/models", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.Models{Object: "list", Data: nil})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)
	_, err := client.Models().ListModels()
	t.Helper()
	if err != nil {
		t.Error(err, "ListModels error")
	}
}

func TestGetModel(t *testing.T) {
	testModelID := "testModelID"
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/models/testModelID", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.Model{Object: "model", ID: testModelID})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)
	model, err := client.Models().GetModel(testModelID)
	t.Helper()
	if err != nil {
		t.Error(err, "GetModel error")
	}
	if model.ID != testModelID {
		t.Errorf("ModelsEndpoint GetModel Model ID mismatch. Got %s. Expected %s", testModelID, model.ID)
	}
}
