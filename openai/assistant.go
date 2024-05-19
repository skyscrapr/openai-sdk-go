package openai

import (
	"fmt"
	"net/url"
	"strconv"
)

const AssistantsEndpointPath = "/assistants/"

// Assistants Endpoint
//
// Manage fine-tuning jobs to tailor a model to your specific training data.
// Related guide: [Fine-tuning models]: https://platform.openai.com/docs/guides/fine-tuning
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/fine-tuning
type AssistantsEndpoint struct {
	*betaEndpoint
}

// Assistants Endpoint
func (c *Client) Assistants() *AssistantsEndpoint {
	return &AssistantsEndpoint{newBetaEndpoint(c, AssistantsEndpointPath)}
}

type Assistant struct {
	Id            string                  `json:"id"`
	Object        string                  `json:"object"` // The object type, which is always assistant.
	CreatedAt     int64                   `json:"created_at"`
	Name          *string                 `json:"name"`
	Description   *string                 `json:"description"`
	Model         string                  `json:"model"`
	Instructions  *string                 `json:"instructions"`
	Tools         []AssistantTool         `json:"tools,omitempty"`
	ToolResources *AssistantToolResources `json:"tool_resources,omitempty"`
	MetaData      map[string]string       `json:"metadata,omitempty"`
	Temperature   float64                 `json:"temperature,omitempty"`
	TopP          float64                 `json:"top_p,omitempty"`
}

type Assistants struct {
	Object  string      `json:"object"`
	Data    []Assistant `json:"data"`
	HasMore bool        `json:"has_more"`
}

type AssistantRequest struct {
	// ID of the model to use. You can use the List models API to see all of your available models, or see our Model overview for descriptions of them.
	Model string `json:"model" binding:"required"`

	// The name of the assistant. The maximum length is 256 characters.
	Name *string `json:"name"`

	// The description of the assistant. The maximum length is 512 characters.
	Description *string `json:"description"`

	// The system instructions that the assistant uses. The maximum length is 32768 characters.
	Instructions *string `json:"instructions"`

	// A list of tool enabled on the assistant. There can be a maximum of 128 tools per assistant. Tools can be of types code_interpreter, retrieval, or function.
	// Defaults to []
	Tools []AssistantTool `json:"tools,omitempty"`

	// A list of file IDs attached to this assistant. There can be a maximum of 20 files attached to the assistant. Files are ordered by their creation date in ascending order.
	ToolResources *AssistantToolResources `json:"tool_resources,omitempty"`

	// Set of 16 key-value pairs that can be attached to an object. This can be useful for storing additional information about the object in a structured format. Keys can be a maximum of 64 characters long and va
	MetaData map[string]string `json:"metadata,omitempty"`

	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
}

type AssistantFile struct {
	Id          string `json:"id"`
	Object      string `json:"object"` // The object type, which is always assistant.file.
	CreatedAt   int64  `json:"created_at"`
	AssistantId string `json:"assistant_id"`
}

type AssistantFiles struct {
	Object  string          `json:"object"`
	Data    []AssistantFile `json:"data"`
	HasMore bool            `json:"has_more"`
}

type AssistantFileRequest struct {
	// A File ID (with purpose="assistants") that the assistant should use. Useful for tools like retrieval and code_interpreter that can access files.
	FileId *string `json:"file_id"`
}

type AssistantTool struct {
	Type     string `json:"type"`
	Function *struct {
		Description *string                `json:"description,omitempty"`
		Name        string                 `json:"name"`
		Parameters  map[string]interface{} `json:"parameters"`
	} `json:"function,omitempty"`
}

type AssistantToolResources struct {
	CodeInterpreter *struct {
		FileIDs []string `json:"file_ids"`
	} `json:"code_interpreter,omitempty"`
	FileSearch *struct {
		VectorStoreIDs []string `json:"vector_store_ids"`
		VectorStores   *struct {
			FileIDs  []string          `json:"file_ids"`
			MetaData map[string]string `json:"metadata,omitempty"`
		} `json:"vector_stores,omitempty"`
	} `json:"file_search,omitempty"`
}

// Create an assistant with a model and instructions.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/assistants/createAssistant
func (e *AssistantsEndpoint) CreateAssistant(req *AssistantRequest) (*Assistant, error) {
	var assistant Assistant
	err := e.do(e, "POST", "", req, nil, &assistant)
	return &assistant, err
}

// Retrieves an assistant.
func (e *AssistantsEndpoint) RetrieveAssistant(assistantId string) (*Assistant, error) {
	var assistant Assistant
	err := e.do(e, "GET", assistantId, nil, nil, &assistant)
	return &assistant, err
}

// Modifies an assistant.
func (e *AssistantsEndpoint) ModifyAssistant(req *AssistantRequest) (*Assistant, error) {
	var assistant Assistant
	err := e.do(e, "POST", "", req, nil, &assistant)
	return &assistant, err
}

// Deletes an assistant.
func (e *AssistantsEndpoint) DeleteAssistant(assistantId string) (bool, error) {
	type DeleteResponse struct {
		Id      string `json:"id"`
		Object  string `json:"object"`
		Deleted bool   `json:"deleted"`
	}
	var resp DeleteResponse
	err := e.do(e, "DELETE", url.QueryEscape(assistantId), nil, nil, &resp)
	if err != nil {
		return false, err
	}
	return resp.Deleted, nil
}

// Returns a list of assistants.
func (e *AssistantsEndpoint) ListAssistants(after *string, limit *int) ([]Assistant, error) {
	v := url.Values{}
	if after != nil {
		v.Add("after", *after)
	}
	if limit != nil {
		v.Add("limit", strconv.Itoa(*limit))
	}
	var assistants Assistants
	err := e.do(e, "GET", "", nil, v, &assistants)
	// TODO: This needs to move somewhere central
	if err == nil && assistants.Object != "list" {
		err = fmt.Errorf("expected 'list' object type, got %s", assistants.Object)
	}
	return assistants.Data, err
}

// Creates an assistant file.
func (e *AssistantsEndpoint) CreateAssistantFile(assistantId string, fileId string) (*AssistantFile, error) {
	req := AssistantFileRequest{
		FileId: &fileId,
	}
	var file AssistantFile
	err := e.do(e, "POST", assistantId, req, nil, &file)
	return &file, err
}

// Retrieves an assistant file.
func (e *AssistantsEndpoint) RetrieveAssistantFile(assistantId string, fileId string) (*AssistantFile, error) {
	var file AssistantFile
	err := e.do(e, "GET", assistantId+"/files/"+fileId, nil, nil, &file)
	return &file, err
}

// Deletes an assistant file.
func (e *AssistantsEndpoint) DeleteAssistantFile(assistantId string, fileId string) (bool, error) {
	var deleted bool
	err := e.do(e, "DELETE", assistantId+"/files/"+fileId, nil, nil, &deleted)
	return deleted, err
}

// Returns a list of assistant files.
func (e *AssistantsEndpoint) ListAssistantFiles(after *string, limit *int) ([]AssistantFile, error) {
	v := url.Values{}
	if after != nil {
		v.Add("after", *after)
	}
	if limit != nil {
		v.Add("limit", strconv.Itoa(*limit))
	}
	var files AssistantFiles
	err := e.do(e, "GET", "", nil, v, &files)
	// TODO: This needs to move somewhere central
	if err == nil && files.Object != "list" {
		err = fmt.Errorf("expected 'list' object type, got %s", files.Object)
	}
	return files.Data, err
}
