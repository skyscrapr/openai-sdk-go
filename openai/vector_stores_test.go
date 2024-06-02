package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skyscrapr/openai-sdk-go/openai"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
)

func TestListVectorStores(t *testing.T) {
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/vector_stores", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.VectorStores{
			Object: "list",
			Data: []openai.VectorStore{
				{
					Id:         "testVectorStoreId",
					Object:     "vector_store",
					CreatedAt:  1613677385,
					Name:       "testVectorStore",
					UsageBytes: 1234,
					FileCounts: &openai.FileCounts{
						InProgress: 1,
						Completed:  1,
						Failed:     0,
						Cancelled:  0,
						Total:      1,
					},
					Status: "testStatus",
					ExpiresAfter: openai.ExpiresAfter{
						Anchor: "last_active_at",
						Days:   1,
					},
					ExpiresAt:    1613677385,
					LastActiveAt: 1613677385,
				},
			},
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)
	_, err := client.VectorStores().ListVectorStores()
	t.Helper()
	if err != nil {
		t.Error(err, "ListVectorStores error")
		t.Fail()
	}
}

func TestCreateVectorStore(t *testing.T) {
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/vector_stores/", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.VectorStore{
			Id:         "testVectorStoreId",
			Object:     "vector_store",
			CreatedAt:  1613677385,
			Name:       "testVectorStore",
			UsageBytes: 1234,
			FileCounts: &openai.FileCounts{
				InProgress: 1,
				Completed:  1,
				Failed:     0,
				Cancelled:  0,
				Total:      1,
			},
			Status: "testStatus",
			ExpiresAfter: openai.ExpiresAfter{
				Anchor: "last_active_at",
				Days:   1,
			},
			ExpiresAt:    1613677385,
			LastActiveAt: 1613677385,
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)

	vectorStore, err := client.VectorStores().CreateVectorStore(&openai.CreateVectorStoresRequest{
		FileIDs: []string{"fileid1"},
		Name:    "testVectorStore",
	})
	t.Helper()
	if err != nil {
		t.Error(err, "RetrieveVectorStore error")
		t.Fail()
	}
	if vectorStore.Id != "testVectorStoreId" {
		t.Error("RetrieveVectorStore ID error")
		t.Fail()
	}
}

func TestRetrieveVectorStores(t *testing.T) {
	testVectorStoreId := "testId"
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/vector_stores/"+testVectorStoreId, func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.VectorStore{
			Id:         "testVectorStoreId",
			Object:     "vector_store",
			CreatedAt:  1613677385,
			Name:       "testVectorStore",
			UsageBytes: 1234,
			FileCounts: &openai.FileCounts{
				InProgress: 1,
				Completed:  1,
				Failed:     0,
				Cancelled:  0,
				Total:      1,
			},
			Status: "testStatus",
			ExpiresAfter: openai.ExpiresAfter{
				Anchor: "last_active_at",
				Days:   1,
			},
			ExpiresAt:    1613677385,
			LastActiveAt: 1613677385,
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)
	_, err := client.VectorStores().RetrieveVectorStore(testVectorStoreId)
	t.Helper()
	if err != nil {
		t.Error(err, "RetrieveVectorStore error")
		t.Fail()
	}
}

func TestDeleteVectorStore(t *testing.T) {
	testVectorStoreId := "testVectorStoreId"
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/vector_stores/"+testVectorStoreId, func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.DeletionStatus{
			Id:      testVectorStoreId,
			Object:  "vector_store.deleted",
			Deleted: true,
		})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)
	_, err := client.VectorStores().DeleteVectorStore(testVectorStoreId)
	t.Helper()
	if err != nil {
		t.Error(err, "DeleteVectorStore error")
		t.Fail()
	}
}
