package openai

import (
	"fmt"
)

const modelsEndpoint = "/models/"

// ModelsEndpoint - OpenAI Models Endpoint
//
//	List and describe the various models available in the API.
//	You can refer to the Models documentation to understand what models are available and the differences between them.
type ModelsEndpoint struct {
	*endpoint
}

// Models - Models Endpoint
func (c *Client) Models() *ModelsEndpoint {
	return &ModelsEndpoint{newEndpoint(c, modelsEndpoint)}
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

// ListModels
//
//	Lists the currently available models, and provides basic information about each one such as the owner and availability.
func (e *ModelsEndpoint) ListModels() ([]Model, error) {
	var models Models
	err := e.do(e, "GET", "", nil, &models)
	if err != nil {
		return nil, err
	}
	if models.Object != "list" {
		return nil, fmt.Errorf("expected 'list' object type, got %s", models.Object)
	}
	return models.Data, nil
}

// GetModel
//
//	Gets the model by ID
func (e *ModelsEndpoint) GetModel(id string) (*Model, error) {
	var model Model
	err := e.do(e, "GET", id, nil, &model)
	if err != nil {
		return nil, err
	}
	return &model, nil
}
