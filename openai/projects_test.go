package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/skyscrapr/openai-sdk-go/openai"
	"github.com/skyscrapr/openai-sdk-go/openai/test"
)

// TestListProjects Tests the Projects endpoint of the API using the mocked server.
func TestListProjects(t *testing.T) {
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/organization/projects", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.Projects{Object: "list", Data: []openai.Project{{Name: "Project1"}}})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)
	_, err := client.Projects().ListProjects()
	t.Helper()
	if err != nil {
		t.Error(err, "TestListProjects error")
	}
}

func TestListProjectsInvalidObject(t *testing.T) {
	expectedError := "expected 'list' object type, got project"

	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/organization/projects", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.Projects{Object: "project", Data: nil})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)
	_, err := client.Projects().ListProjects()
	t.Helper()
	if err != nil && err.Error() != expectedError {
		t.Errorf("Unexpected error: %v , expected: %s", err, expectedError)
		t.Fail()
	}
}

func TestRetrieveProject(t *testing.T) {
	testProjectID := "testProjectID"
	ts := openai_test.NewTestServer()
	ts.RegisterHandler("/v1/organization/projects/testProjectID", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(openai.Project{Object: "project", ID: testProjectID})
		fmt.Fprintln(w, string(resBytes))
	})
	ts.HTTPServer.Start()
	defer ts.HTTPServer.Close()

	client := openai_test.NewTestClient(ts)
	project, err := client.Projects().RetrieveProject(testProjectID)
	t.Helper()
	if err != nil {
		t.Error(err, "GetProject error")
	}
	if project.ID != testProjectID {
		t.Errorf("ProjectsEndpoint GetProject Project ID mismatch. Got %s. Expected %s", testProjectID, project.ID)
	}
}
