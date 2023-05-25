package openai

import (
	"encoding/json"
	"fmt"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
	"net/http"
	"testing"
)

var (
	testModels = []Model{}
)

func TestListModels(t *testing.T) {
	// testModels := []Model{
	// 	{
	// 		CreatedAt:  0,
	// 		ID:         "",
	// 		Object:     "",
	// 		OwnedBy:    "",
	// 		Permission: []Permission{},
	// 		Root:       "",
	// 		Parent:     "",
	// 	},
	// }
	// testServer := test.NewTestServer()
	// //NewTestServer(t, "GET", "/models", testModels)
	// defer testServer.Close()
	// testClient := newTestClient(t, testServer)
	// e := testClient.Models()
	// models, err := e.ListModels()
	// if err != nil {
	// 	t.Errorf("Unexpected Errro: %s", err)
	// }
	// if models == nil {
	// 	t.Fail()
	// }
	// testCheckStructEqual(t, models, testModels)

	server := test.NewTestServer()
	server.RegisterHandler("/v1/models", handlerListModels)

	// create the test server
	testServer := server.OpenAITestServer()
	testServer.Start()
	defer testServer.Close()

	testClient := newTestClient(t, testServer)
	models, err := testClient.Models().ListModels()
	if err != nil {
		t.Errorf("Unexpected Errro: %s", err)
	}
	if models == nil {
		t.Fail()
	}
	test.CheckStructEqual(t, models, testModels)
}

// handleModelsEndpoint Handles the models endpoint by the test server.
func handlerListModels(w http.ResponseWriter, _ *http.Request) {
	resBytes, _ := json.Marshal([]Model{})
	fmt.Fprintln(w, string(resBytes))
}
