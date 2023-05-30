package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skyscrapr/openai-sdk-go/openai"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
)

func TestListFiles(t *testing.T) {
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/files", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.Files{
			Object: "list",
			Data: []openai.File{
				{
					Id:        "testFileId",
					Object:    "file",
					Bytes:     175,
					CreatedAt: 1613677385,
					Filename:  "train.jsonl",
					//Purpose: FineTune{},
				},
			},
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)
	_, err := client.Files().ListFiles()
	t.Helper()
	if err != nil {
		t.Error(err, "ListFiles error")
		t.Fail()
	}
}

func TestUploadFile(t *testing.T) {
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/files", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.File{
			Id:        "testFileId",
			Object:    "file",
			Bytes:     175,
			CreatedAt: 1613677385,
			Filename:  "train.jsonl",
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)
	_, err := client.Files().UploadFile(
		&openai.UploadFileRequest{
			File:    "testFileName",
			Purpose: "fine-tune",
		},
	)
	t.Helper()
	if err != nil {
		t.Error(err, "UploadFile error")
		t.Fail()
	}
}

func TestDeleteFile(t *testing.T) {
	testFileId := "testFileId"
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/files/"+testFileId, func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.DeleteFileResponse{
			Id:      testFileId,
			Object:  "file",
			Deleted: true,
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)
	_, err := client.Files().DeleteFile(testFileId)
	t.Helper()
	if err != nil {
		t.Error(err, "DeleteFile error")
		t.Fail()
	}
}

func TestRetrieveFile(t *testing.T) {
	testFileId := "testFileId"
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/files/"+testFileId, func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.File{
			Id:        "testFileId",
			Object:    "file",
			Bytes:     175,
			CreatedAt: 1613677385,
			Filename:  "train.jsonl",
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)
	_, err := client.Files().RetrieveFile(testFileId)
	t.Helper()
	if err != nil {
		t.Error(err, "DeleteFile error")
		t.Fail()
	}
}

// TODO: Add test for RetrieveFileContent
