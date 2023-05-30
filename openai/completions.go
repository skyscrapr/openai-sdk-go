package openai

const CompletionsEndpointPath = "/completions/"

// Completions Endpoint
//
//	Given a prompt, the model will return one or more predicted completions,
//  and can also return the probabilities of alternative tokens at each position.
type CompletionsEndpoint struct {
	*endpoint
}

// Completions Endpoint
func (c *Client) Completions() *CompletionsEndpoint {
	return &CompletionsEndpoint{newEndpoint(c, CompletionsEndpointPath)}
}

type CompletionRequest struct {
	// ID of the model to use. 
	// You can use the [List models]: https://platform.openai.com/docs/api-reference/models/list API to see all of your available models,
	// or see our [Model overview]: https://platform.openai.com/docs/models/overview for descriptions of them.
	Model string `json:"model"`
	// The prompt(s) to generate completions for, encoded as a string, array of strings, array of tokens, or array of token arrays.
	// Note that <|endoftext|> is the document separator that the model sees during training,
	// so if a prompt is not specified the model will generate as if from the beginning of a new document.
	Prompt []string `json:"prompt,omitempty"`
	// MaxTokens int `json:"max_tokens,omitempty" binding:"omitempty,max=4096"`
	// Temperature float32 `json:"temperature,omitempty"`
	// N int `json:"n,omitempty"`
	// Stop []string `json:"stop,omitempty"`
}

type CompletionResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}



// Creates a completion for the provided prompt and parameters.
//
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/completions/create
func (e *CompletionsEndpoint) CreateCompletion(req *CompletionRequest) (*CompletionResponse, error) {
	var resp CompletionResponse
	err := e.do(e, "POST", "", req, &resp)
	return &resp, err
}
