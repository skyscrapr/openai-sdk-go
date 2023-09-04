package openai

const EditsEndpointPath = "/edits/"

// Edits Endpoint
//
//	Given a prompt and an instruction, the model will return an edited version of the prompt.
type EditsEndpoint struct {
	*endpoint
}

// Edits Endpoint
func (c *Client) Edits() *EditsEndpoint {
	return &EditsEndpoint{newEndpoint(c, EditsEndpointPath)}
}

type EditRequest struct {
	// ID of the model to use. You can use the text-davinci-edit-001 or code-davinci-edit-001 model with this endpoint.
	Model string `json:"model" binding:"required"`
	// Defaults to ''
	// The input text to use as a starting point for the edit.
	Input string `json:"input,omitempty"`
	// The instruction that tells the model how to edit the prompt.
	Instruction string `json:"instruction" binding:"required"`
	// Defaults to 1
	// How many edits to generate for the input and instruction.
	N int `json:"n,omitempty"`
	// Defaults to 1
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic.
	// We generally recommend altering this or top_p but not both.
	Temperature int `json:"temperature,omitempty"`
	// Defaults to 1
	// An alternative to sampling with temperature, called nucleus sampling, where the model considers the results of the tokens with top_p probability mass. So 0.1 means only the tokens comprising the top 10% probability mass are considered.
	// We generally recommend altering this or temperature but not both.
	TopP int `json:"top_p,omitempty"`
}

type EditResponse struct {
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Text  string `json:"text"`
		Index int    `json:"index"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Creates a new edit for the provided input, instruction, and parameters.
//
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/edits/create
func (e *EditsEndpoint) CreateEdit(req *EditRequest) (*EditResponse, error) {
	var resp EditResponse
	err := e.do(e, "POST", "", req, nil, &resp)
	return &resp, err
}
