package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skyscrapr/openai-sdk-go/openai"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
)

func TestCreateFineTune(t *testing.T) {
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/audio/transcriptions", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.AudioResponse{
			Text: "test",
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)

	req := openai.AudioTranscriptionRequest{
		Model: "test",
	}
	_, err := client.Audio().CreateTranscription(&req)
	t.Helper()
	if err != nil {
		t.Error(err, "CreateTranscription error")
		t.Fail()
	}
}
