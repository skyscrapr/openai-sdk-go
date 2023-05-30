package openai

const AudioEndpointPath = "/audio/"

// Audio Endpoint
//
//		Learn how to turn audio into text.
//	 Related guide: [Speech to text]: https://platform.openai.com/docs/guides/speech-to-text
type AudioEndpoint struct {
	*endpoint
}

// Audio Endpoint
func (c *Client) Audio() *AudioEndpoint {
	return &AudioEndpoint{newEndpoint(c, AudioEndpointPath)}
}

type AudioResponse struct {
	Text string `json:"text"`
}

type AudioTranscriptionRequest struct {
	// The audio file to transcribe, in one of these formats: mp3, mp4, mpeg, mpga, m4a, wav, or webm.
	File string `json:"file" binding:"required"`
	// 	ID of the model to use. Only whisper-1 is currently available.
	Model string `json:"model" binding:"required"`
	// An optional text to guide the model's style or continue a previous audio segment. The prompt should match the audio language.
	Prompt string `json:"prompt,omitempty"`
	// 	Defaults to json
	// The format of the transcript output, in one of these options: json, text, srt, verbose_json, or vtt.
	ResponseFormat string `json:"response_format,omitempty"`
	// Defaults to 0
	// The sampling temperature, between 0 and 1. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic. If set to 0, the model will use log probability to automatically increase the temperature until certain thresholds are hit.
	Temperature int `json:"temperature,omitempty"`
	// The language of the input audio. Supplying the input language in ISO-639-1 format will improve accuracy and latency.
	Language string `json:"language,omitempty"`
}

type AudioTranslationRequest struct {
	// The audio file to translate, in one of these formats: mp3, mp4, mpeg, mpga, m4a, wav, or webm.
	File string `json:"file" binding:"required"`
	// 	ID of the model to use. Only whisper-1 is currently available.
	Model string `json:"model" binding:"required"`
	// An optional text to guide the model's style or continue a previous audio segment. The prompt should be in English.
	Prompt string `json:"prompt,omitempty"`
	// 	Defaults to json
	// The format of the transcript output, in one of these options: json, text, srt, verbose_json, or vtt.
	ResponseFormat string `json:"response_format,omitempty"`
	// Defaults to 0
	// The sampling temperature, between 0 and 1. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic. If set to 0, the model will use log probability to automatically increase the temperature until certain thresholds are hit.
	Temperature int `json:"temperature,omitempty"`
}

// Transcribes audio into the input language.
//
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/audio/create
func (e *AudioEndpoint) CreateTranscription(req *AudioTranscriptionRequest) (*AudioResponse, error) {
	var resp AudioResponse
	err := e.do(e, "POST", "transcriptions", req, &resp)
	return &resp, err
}

// Translates audio into into English.
//
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/audio/create
func (e *AudioEndpoint) CreateTranslation(req *AudioTranslationRequest) (*AudioResponse, error) {
	var resp AudioResponse
	err := e.do(e, "POST", "translations", req, &resp)
	return &resp, err
}
