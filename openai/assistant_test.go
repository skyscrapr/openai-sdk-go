package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skyscrapr/openai-sdk-go/openai"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
)

func TestCreateAssistant(t *testing.T) {
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/assistants", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.FineTuningJob{
			Id:     "test_id",
			Object: "assistant",
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)

	name := "test_name"
	description := "test_description"
	instructions := "test_instructions"
	req := openai.AssistantRequest{
		Model:        "test_model",
		Name:         &name,
		Description:  &description,
		Instructions: &instructions,
		FileIds:      []string{"test_file_id_1", "test_file_id_2"},
		MetaData:     map[string]string{"test_key_1": "test_value_1"},
	}
	_, err := client.Assistants().CreateAssistant(&req)
	t.Helper()
	if err != nil {
		t.Error(err, "CreateAssistant error")
		t.Fail()
	}
}
