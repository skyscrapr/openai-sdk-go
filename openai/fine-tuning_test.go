package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skyscrapr/openai-sdk-go/openai"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
)

func TestCreateFineTuningJob(t *testing.T) {
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/fine_tuning/jobs", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.FineTuningJob{
			Id:     "test_id",
			Object: "fine_tuning.job",
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)

	req := openai.CreateFineTuningJobRequest{
		TrainingFile:   "file_id_1",
		ValidationFile: "file_id_2",
		Model:          "test_model",
		Suffix:         "test_suffix",
	}
	_, err := client.FineTuning().CreateFineTuningJob(&req)
	t.Helper()
	if err != nil {
		t.Error(err, "CreateFineTuningJob error")
		t.Fail()
	}
}
