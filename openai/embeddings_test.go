package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skyscrapr/openai-sdk-go/openai"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
)

func TestCreateEmbeddings(t *testing.T) {
	testModelID := "testModelID"
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/embeddings", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.EmbeddingsResponse{
			Model: testModelID,
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)

	req := openai.EmbeddingsRequest{
		Model: testModelID,
	}
	resp, err := client.Embeddings().CreateEmbeddings(&req)
	t.Helper()
	if err != nil {
		t.Error(err, "CreateEmbeddings error")
		t.Fail()
	}
	if resp.Model != testModelID {
		t.Errorf("Embeddings Endpoint CreateEmbeddings Model ID mismatch. Got %s. Expected %s", resp.Model, testModelID)
	}
}
