package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skyscrapr/openai-sdk-go/openai"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
)

func TestCreateEdit(t *testing.T) {
	testModelID := "testModelID"
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/edits", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.EditResponse{
			Created: 1,
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)

	req := openai.EditRequest{
		Model:       testModelID,
		Instruction: "test",
	}
	_, err := client.Edits().CreateEdit(&req)
	t.Helper()
	if err != nil {
		t.Error(err, "CreateEdit error")
		t.Fail()
	}
}
