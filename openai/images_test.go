package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skyscrapr/openai-sdk-go/openai"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
)

func TestCreateImage(t *testing.T) {
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/images/generations", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.ImagesResponse{
			Created: 1,
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)

	req := openai.CreateImageRequest{
		Prompt: "test",
	}
	_, err := client.Images().CreateImage(&req)
	t.Helper()
	if err != nil {
		t.Error(err, "CreateImage error")
		t.Fail()
	}
}

func TestCreateImageEdit(t *testing.T) {
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/images/edits", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.ImagesResponse{
			Created: 1,
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)

	req := openai.CreateImageEditRequest{
		Image:  "test",
		Prompt: "test",
	}
	_, err := client.Images().CreateImageEdit(&req)
	t.Helper()
	if err != nil {
		t.Error(err, "CreateImageEdit error")
		t.Fail()
	}
}

func TestCreateImageVariation(t *testing.T) {
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/images/edits", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.ImagesResponse{
			Created: 1,
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)

	req := openai.CreateImageVariationRequest{
		Image: "test",
	}
	_, err := client.Images().CreateImageVariation(&req)
	t.Helper()
	if err != nil {
		t.Error(err, "CreateImageVariation error")
		t.Fail()
	}
}
