package openai

import (
	"fmt"
	"net/url"
)

const ProjectsEndpointPath = "/organization/projects"

// ProjectsEndpoint - OpenAI Projects Endpoint
//
//	List and describe the projects available.
//	You can refer to the [Projects]: https://platform.openai.com/docs/api-reference/projects documentation.
type ProjectsEndpoint struct {
	*organizationEndpoint
}

// Projects - Projects Endpoint
func (c *Client) Projects() *ProjectsEndpoint {
	return &ProjectsEndpoint{newOrganizationEndpoint(c, ProjectsEndpointPath)}
}

type Project struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	Name       string `json:"name"`
	CreatedAt  int64  `json:"created_at"`
	ArchivedAt int64  `json:"archived_at"`
	Status     string `json:"status"`
}

type Projects struct {
	Object string    `json:"object"`
	Data   []Project `json:"data"`
}

type ProjectRequest struct {
	// The name of the assistant. The maximum length is 256 characters.
	Name *string `json:"name"`
}

type ProjectServiceAccount struct {
	ID        string         `json:"id"`
	ProjectID string         `json:"-"`
	Object    string         `json:"object"`
	Name      string         `json:"name"`
	Role      string         `json:"role"`
	CreatedAt int64          `json:"created_at"`
	ApiKey    *ProjectApiKey `json:"api_key,omitempty"`
}

type ProjectApiKey struct {
	ID        string  `json:"id"`
	Object    string  `json:"object"`
	Value     string  `json:"value"`
	Name      *string `json:"name"`
	CreatedAt int64   `json:"created_at"`
}

type ProjectServiceAccounts struct {
	Object string                  `json:"object"`
	Data   []ProjectServiceAccount `json:"data"`
}

type ProjectServiceAccountRequest struct {
	// The name of the service account being created.
	Name *string `json:"name"`
}

// Lists the currently available projects,
// and provides basic information about each one.
//
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/projects/list
func (e *ProjectsEndpoint) ListProjects() ([]Project, error) {
	var projects Projects
	err := e.do(e, "GET", "", nil, nil, &projects)
	// TODO: This needs to move somewhere central
	if err == nil && projects.Object != "list" {
		err = fmt.Errorf("expected 'list' object type, got %s", projects.Object)
	}
	return projects.Data, err
}

// Create a project.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/projects/create
func (e *ProjectsEndpoint) CreateProject(req *ProjectRequest) (*Project, error) {
	var project Project
	err := e.do(e, "POST", "", req, nil, &project)
	return &project, err
}

// Retrieves a project instance,
// providing basic information about the project.
//
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/projects/retrieve
func (e *ProjectsEndpoint) RetrieveProject(id string) (*Project, error) {
	var project Project
	err := e.do(e, "GET", id, nil, nil, &project)
	return &project, err
}

// Modifies a project instance,
//
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/projects/modify
func (e *ProjectsEndpoint) ModifyProject(id string, req ProjectRequest) (*Project, error) {
	var project Project
	err := e.do(e, "POST", id, req, nil, &project)
	return &project, err
}

// Archive a project.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/projects/archive
func (e *ProjectsEndpoint) ArchiveProject(id string) (*Project, error) {
	var project Project
	err := e.do(e, "POST", url.QueryEscape(id)+"/archive", nil, nil, &project)
	return &project, err
}

// Returns a list of service accounts in the project.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/project-service-accounts/list
func (e *ProjectsEndpoint) ListProjectServiceAccounts(projectId string) ([]ProjectServiceAccount, error) {
	var projectServiceAccounts ProjectServiceAccounts
	err := e.do(e, "GET", url.QueryEscape(projectId)+"/service_accounts", nil, nil, &projectServiceAccounts)
	if err == nil && projectServiceAccounts.Object != "list" {
		err = fmt.Errorf("expected 'list' object type, got %s", projectServiceAccounts.Object)
	}
	return projectServiceAccounts.Data, err
}

// Creates a new service account in the project. This also returns an unredacted API key for the service account.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/project-service-accounts/create
func (e *ProjectsEndpoint) CreateProjectServiceAccount(projectId string, req *ProjectServiceAccountRequest) (*ProjectServiceAccount, error) {
	var projectServiceAccount ProjectServiceAccount
	err := e.do(e, "POST", url.QueryEscape(projectId)+"/service_accounts", req, nil, &projectServiceAccount)
	return &projectServiceAccount, err
}

// Retrieves a service account in the project.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/project-service-accounts/retrieve
func (e *ProjectsEndpoint) RetrieveProjectServiceAccount(projectId string, id string) (*ProjectServiceAccount, error) {
	var projectServiceAccount ProjectServiceAccount
	err := e.do(e, "GET", url.QueryEscape(projectId)+"/service_accounts/"+id, nil, nil, &projectServiceAccount)
	return &projectServiceAccount, err
}

// Deletes a service account from the project.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/project-service-accounts/delete
func (e *ProjectsEndpoint) DeleteProjectServiceAccount(projectId string, id string) (bool, error) {
	type DeleteResponse struct {
		Id      string `json:"id"`
		Object  string `json:"object"`
		Deleted bool   `json:"deleted"`
	}
	var resp DeleteResponse
	err := e.do(e, "DELETE", url.QueryEscape(projectId)+"/service_accounts/"+id, nil, nil, &resp)
	if err != nil {
		return false, err
	}
	return resp.Deleted, nil
}
