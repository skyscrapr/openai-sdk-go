package openai

import (
	"fmt"
	"net/url"
)

const ModelsEndpointPath = "/models/"

// ModelsEndpoint - OpenAI Models Endpoint
//
//	List and describe the various models available in the API.
//	You can refer to the [Models]: https://platform.openai.com/docs/models documentation to understand what models are available and the differences between them.
type ModelsEndpoint struct {
	*endpoint
}

// Models - Models Endpoint
func (c *Client) Models() *ModelsEndpoint {
	return &ModelsEndpoint{newEndpoint(c, ModelsEndpointPath)}
}

// Model - OpenAPI Model.
type Model struct {
	CreatedAt  int64        `json:"created"`
	ID         string       `json:"id"`
	Object     string       `json:"object"`
	OwnedBy    string       `json:"owned_by"`
	Permission []Permission `json:"permission"`
	Root       string       `json:"root"`
	Parent     string       `json:"parent"`
}

// Permission - OpenAPI Permission.
type Permission struct {
	CreatedAt          int64       `json:"created"`
	ID                 string      `json:"id"`
	Object             string      `json:"object"`
	AllowCreateEngine  bool        `json:"allow_create_engine"`
	AllowSampling      bool        `json:"allow_sampling"`
	AllowLogprobs      bool        `json:"allow_logprobs"`
	AllowSearchIndices bool        `json:"allow_search_indices"`
	AllowView          bool        `json:"allow_view"`
	AllowFineTuning    bool        `json:"allow_fine_tuning"`
	Organization       string      `json:"organization"`
	Group              interface{} `json:"group"`
	IsBlocking         bool        `json:"is_blocking"`
}

type Models struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// Lists the currently available models,
// and provides basic information about each one such as the owner and availability.
//
// [OpenAI Documentation]: https://beta.openai.com/docs/api-reference/models/list
func (e *ModelsEndpoint) ListModels() ([]Model, error) {
	var models Models
	err := e.do(e, "GET", "", nil, nil, &models)
	// TODO: This needs to move somewhere central
	if err == nil && models.Object != "list" {
		err = fmt.Errorf("expected 'list' object type, got %s", models.Object)
	}
	return models.Data, err
}

// Retrieves a model instance,
// providing basic information about the model such as the owner and permissioning.
//
// [OpenAI Documentation]: https://beta.openai.com/docs/api-reference/models/retrieve
func (e *ModelsEndpoint) RetrieveModel(id string) (*Model, error) {
	var model Model
	err := e.do(e, "GET", id, nil, nil, &model)
	return &model, err
}

// Delete a fine-tuned model. You must have the Owner role in your organization.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/fine-tunes/delete-model
func (e *ModelsEndpoint) DeleteFineTuneModel(id string) (bool, error) {
	type DeleteResponse struct {
		Id      string `json:"id"`
		Object  string `json:"object"`
		Deleted bool   `json:"deleted"`
	}
	var resp DeleteResponse
	err := e.do(e, "DELETE", url.QueryEscape(id), nil, nil, &resp)
	if err != nil {
		return false, err
	}
	return resp.Deleted, nil
}
