package openai

const ModerationsEndpointPath = "/moderations/"

// Moderations Endpoint
//
//	Given a input text, outputs if the model classifies it as violating OpenAI's content policy.
//	Related guide: [Moderations]: https://platform.openai.com/docs/guides/moderation
type ModerationsEndpoint struct {
	*endpoint
}

// Moderations Endpoint
func (c *Client) Moderations() *ModerationsEndpoint {
	return &ModerationsEndpoint{newEndpoint(c, ModerationsEndpointPath)}
}

type ModerationRequest struct {
	// The input text to classify
	Input []string `json:"input" binding:"required"`
	// Defaults to text-moderation-latest
	// Two content moderations models are available: text-moderation-stable and text-moderation-latest.
	// The default is text-moderation-latest which will be automatically upgraded over time. This ensures you are always using our most accurate model. If you use text-moderation-stable, we will provide advanced notice before updating the model. Accuracy of text-moderation-stable may be slightly lower than for text-moderation-latest.Model string `json:"model" binding:"required"`
	Model string `json:"model.omntempty"`
}

type Moderation struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Results []struct {
		Categories []struct {
			Hate            bool `json:"hate"`
			HateThreatening bool `json:"hate/threatening"`
			SelfHarm        bool `json:"self-harm"`
			Sexual          bool `json:"sexual"`
			SexualMinors    bool `json:"sexual/minors"`
			Violence        bool `json:"violence"`
			ViolenceGraphic bool `json:"violence/graphic"`
		} `json:"categories"`
		CategoryScores []struct {
			Hate            float64 `json:"hate"`
			HateThreatening float64 `json:"hate/threatening"`
			SelfHarm        float64 `json:"self-harm"`
			Sexual          float64 `json:"sexual"`
			SexualMinors    float64 `json:"sexual/minors"`
			Violence        float64 `json:"violence"`
			ViolenceGraphic float64 `json:"violence/graphic"`
		} `json:"category_scores"`
		Flagged bool `json:"flagged"`
	} `json:"results"`
}

// Classifies if text violates OpenAI's Content Policy
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/moderations/create
func (e *ModerationsEndpoint) CreateModeration(req *ModerationRequest) (*Moderation, error) {
	var moderation Moderation
	err := e.do(e, "POST", "", req, &moderation)
	return &moderation, err
}
