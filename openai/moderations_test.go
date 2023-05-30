package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skyscrapr/openai-sdk-go/openai"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
)

func TestCreateModeration(t *testing.T) {
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/moderations", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.Moderation{
			Id: "testModerationId",
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)

	req := openai.ModerationRequest{
		Input: []string{"test"},
	}
	_, err := client.Moderations().CreateModeration(&req)
	t.Helper()
	if err != nil {
		t.Error(err, "CreateModerations error")
		t.Fail()
	}
}
