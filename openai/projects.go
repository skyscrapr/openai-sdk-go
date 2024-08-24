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

// Project - OpenAPI Project.
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
