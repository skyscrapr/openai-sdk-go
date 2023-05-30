package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skyscrapr/openai-sdk-go/openai"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
)

func TestCreateCompletion(t *testing.T) {
	testModelID := "testModelID"
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/completions", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.CompletionResponse{
			Model: testModelID,
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)

	req := openai.CompletionRequest{
		Model:     "ada",
	}
	req.Prompt = []string{"Lorem ipsum"}
	resp, err := client.Completions().CreateCompletion(&req)
	t.Helper()
	if err != nil {
		t.Error(err, "CreateCompletion error")
		t.Fail()
	}
	if resp.Model != testModelID {
		t.Errorf("Completions Endpoint CreateCompletion Model ID mismatch. Got %s. Expected %s", resp.Model, testModelID)
	}
}
