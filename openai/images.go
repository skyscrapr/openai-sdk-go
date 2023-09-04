package openai

const ImagesEndpointPath = "/images/"

// Images Endpoint
//
//		Given a prompt and/or an input image, the model will generate a new image.
//	 Related guide: [Image generation]: https://platform.openai.com/docs/guides/images
type ImagesEndpoint struct {
	*endpoint
}

// Images Endpoint
func (c *Client) Images() *ImagesEndpoint {
	return &ImagesEndpoint{newEndpoint(c, ImagesEndpointPath)}
}

type ImagesResponse struct {
	Created int `json:"created"`
	Data    []struct {
		Url string `json:"url"`
	} `json:"data"`
}

type CreateImageRequest struct {
	// A text description of the desired image(s). The maximum length is 1000 characters.
	Prompt string `json:"prompt" binding:"required"`
	// Defaults to 1
	// The number of images to generate. Must be between 1 and 10.
	N int `json:"n,omitempty"`
	// 	Defaults to 1024x1024
	// The size of the generated images. Must be one of 256x256, 512x512, or 1024x1024.
	Size string `json:"size,omitempty"`
	// Defaults to url
	// The format in which the generated images are returned. Must be one of url or b64_json.
	ResponseFormat string `json:"response_format,omitempty"`
	// A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse. Learn more.
	User string `json:"user,omitempty"`
}

type CreateImageEditRequest struct {
	// The image to edit. Must be a valid PNG file, less than 4MB, and square. If mask is not provided, image must have transparency, which will be used as the mask.
	Image string `json:"image" binding:"required"`
	// An additional image whose fully transparent areas (e.g. where alpha is zero) indicate where image should be edited. Must be a valid PNG file, less than 4MB, and have the same dimensions as image.
	Mask string `json:"mask,omitempty"`
	// A text description of the desired image(s). The maximum length is 1000 characters.
	Prompt string `json:"prompt" binding:"required"`
	// Defaults to 1
	// The number of images to generate. Must be between 1 and 10.
	N int `json:"n,omitempty"`
	// Defaults to 1024x1024
	// The size of the generated images. Must be one of 256x256, 512x512, or 1024x1024.
	Size string `json:"size,omitempty"`
	// Defaults to url
	// The format in which the generated images are returned. Must be one of url or b64_json.
	ResponseFormat string `json:"response_format,omitempty"`
	// A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse. Learn more.
	User string `json:"user,omitempty"`
}

type CreateImageVariationRequest struct {
	// The image to use as the basis for the variation(s). Must be a valid PNG file, less than 4MB, and square.
	Image string `json:"image" binding:"required"`
	// Defaults to 1
	// The number of images to generate. Must be between 1 and 10.
	N int `json:"n,omitempty"`
	// Defaults to 1024x1024
	// The size of the generated images. Must be one of 256x256, 512x512, or 1024x1024.
	Size string `json:"size,omitempty"`
	// Defaults to url
	// The format in which the generated images are returned. Must be one of url or b64_json.
	ResponseFormat string `json:"response_format,omitempty"`
	// A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse. Learn more.
	User string `json:"user,omitempty"`
}

// Creates an image given a prompt.
//
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/images/create
func (e *ImagesEndpoint) CreateImage(req *CreateImageRequest) (*ImagesResponse, error) {
	var resp ImagesResponse
	err := e.do(e, "POST", "generations", req, nil, &resp)
	return &resp, err
}

// Creates an edited or extended image given an original image and a prompt.
//
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/images/create-edit
func (e *ImagesEndpoint) CreateImageEdit(req *CreateImageEditRequest) (*ImagesResponse, error) {
	var resp ImagesResponse
	err := e.do(e, "POST", "edits", req, nil, &resp)
	return &resp, err
}

// Creates a variation of a given image.
//
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/images/create-variation
func (e *ImagesEndpoint) CreateImageVariation(req *CreateImageVariationRequest) (*ImagesResponse, error) {
	var resp ImagesResponse
	err := e.do(e, "POST", "edits", req, nil, &resp)
	return &resp, err
}
