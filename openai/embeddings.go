package openai

const EmbeddingsEndpointPath = "/embeddings/"

// Embeddings Endpoint
//
//		Get a vector representation of a given input that can be easily consumed by machine learning models and algorithms.
//	 Related guide: [Embeddings]: https://platform.openai.com/docs/guides/embeddings
type EmbeddingsEndpoint struct {
	*endpoint
}

// Completions Endpoint
func (c *Client) Embeddings() *EmbeddingsEndpoint {
	return &EmbeddingsEndpoint{newEndpoint(c, EmbeddingsEndpointPath)}
}

type EmbeddingsRequest struct {
	// ID of the model to use.
	// You can use the [List models]: https://platform.openai.com/docs/api-reference/models/list API to see all of your available models,
	// or see our [Model overview]: https://platform.openai.com/docs/models/overview for descriptions of them.
	Model string `json:"model" binding:"required"`
	// Input text to get embeddings for, encoded as a string or array of tokens. To get embeddings for multiple inputs in a single request, pass an array of strings or array of token arrays. Each input must not exceed 8192 tokens in length.
	Input string `json:"input" binding:"required"`
	// A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse. Learn more.
	User string `json:"user,omitempty"`
}

type EmbeddingsResponse struct {
	Object string `json:"object"`
	Model  string `json:"model"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float64 `json:"embedding,omitempty"`
	} `json:"data"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// Creates an embedding vector representing the input text.
//
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/embeddings
func (e *EmbeddingsEndpoint) CreateEmbeddings(req *EmbeddingsRequest) (*EmbeddingsResponse, error) {
	var resp EmbeddingsResponse
	err := e.do(e, "POST", "", req, &resp)
	return &resp, err
}
