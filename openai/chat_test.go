package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skyscrapr/openai-sdk-go/openai"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
)

func TestCreateChatCompletion(t *testing.T) {
	testModelID := "testModelID"
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/chat/completions", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.CompletionResponse{
			Model: testModelID,
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)

	req := openai.ChatCompletionRequest{
		Model:    testModelID,
		Messages: []string{"test"},
		Role:     "user",
		Content:  "What is the capital of France?",
	}
	resp, err := client.Chat().CreateChatCompletion(&req)
	t.Helper()
	if err != nil {
		t.Error(err, "CreateChatCompletion error")
		t.Fail()
	}
	if resp.Model != testModelID {
		t.Errorf("Chat Endpoint CreateChatCompletion Model ID mismatch. Got %s. Expected %s", resp.Model, testModelID)
		t.Fail()
	}
}
