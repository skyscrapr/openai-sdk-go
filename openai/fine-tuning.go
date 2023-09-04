package openai

import (
	"fmt"
	"net/url"
	"strconv"
)

const FineTuningEndpointPath = "/fine_tuning/"

// FineTuning Endpoint
//
// Manage fine-tuning jobs to tailor a model to your specific training data.
// Related guide: [Fine-tuning models]: https://platform.openai.com/docs/guides/fine-tuning
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/fine-tuning
type FineTuningEndpoint struct {
	*endpoint
}

// Files Endpoint
func (c *Client) FineTuning() *FineTuningEndpoint {
	return &FineTuningEndpoint{newEndpoint(c, FineTuningEndpointPath)}
}

type FineTuningJob struct {
	Id             string `json:"id"`
	Object         string `json:"object"`
	CreatedAt      int64  `json:"created_at"`
	FinishedAt     int64  `json:"finished_at"`
	Model          string `json:"model"`
	FineTunedModel string `json:"fine_tuned_model"`
	OrganizationId string `json:"organization_id"`
	Status         string `json:"status"`
	Hyperparams    struct {
		NEpochs int64 `json:"n_epochs,omitempty"`
	} `json:"hyperparams,omitempty"`
	TrainingFile   string   `json:"training_file"`
	ValidationFile string   `json:"validation_file"`
	ResultFiles    []string `json:"result_files"`
	TrainedTokens  int64    `json:"trained_tokens,omitempty"`
}

type FineTuningEvents struct {
	Object  string            `json:"object"`
	Data    []FineTuningEvent `json:"data"`
	HasMore bool              `json:"has_more"`
}

type FineTuningEvent struct {
	Object    string `json:"object"`
	Id        string `json:"id"`
	CreatedAt int    `json:"created_at"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	// Data	  string `json:"data"`
	Type string `json:"type"`
}

// EventHandler is a callback that gets called every time event on the SSE
// stream is received. Error returned from handler function will be passed to
// the error handler.
//
// Users of this package have to provide this function implementation.
// type FineTuneEventHandler func(e *FineTuneEvent) error

type CreateFineTuningJobRequest struct {
	// The ID of an uploaded file that contains training data.
	// See [upload file]: https://platform.openai.com/docs/api-reference/files/upload for how to upload a file.
	// Your dataset must be formatted as a JSONL file, where each training example is a JSON object with the keys "prompt" and "completion". Additionally, you must upload your file with the purpose fine-tune.
	// See the [fine-tuning guide]: https://platform.openai.com/docs/guides/fine-tuning/creating-training-data for more details.
	TrainingFile string `json:"training_file" binding:"required"`

	// The ID of an uploaded file that contains validation data.
	// If you provide this file, the data is used to generate validation metrics periodically during fine-tuning.
	// These metrics can be viewed in the [fine-tuning results file]: https://platform.openai.com/docs/guides/fine-tuning/analyzing-your-fine-tuned-model. Your train and validation data should be mutually exclusive.
	// Your dataset must be formatted as a JSONL file, where each validation example is a JSON object with the keys "prompt" and "completion". Additionally, you must upload your file with the purpose fine-tune.
	// See the [fine-tuning guide]: https://platform.openai.com/docs/guides/fine-tuning/creating-training-data for more details.
	ValidationFile string `json:"validation_file,omitempty"`

	// The name of the base model to fine-tune. You can select one of the supported models.
	Model string `json:"model,omitempty" binding:"required"`

	// The hyperparameters used for the fine-tuning job.
	Hyperparameters struct {
		// The number of epochs to train the model for. An epoch refers to one full cycle through the training dataset.
		NEpochs int64 `json:"n_epochs,omitempty"`
	} `json:"hyperparameters,omitempty"`

	// Defaults to null
	// A string of up to 40 characters that will be added to your fine-tuned model name.
	// For example, a suffix of "custom-model-name" would produce a model name like ft:gpt-3.5-turbo:openai:custom-model-name:7p4lURel.
	Suffix string `json:"suffix,omitempty"`
}

// Creates a job that fine-tunes a specified model from a given dataset.
// Response includes details of the enqueued job including job status and the name of the fine-tuned models once complete.
// Learn more about Fine-tuning
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/fine-tunes
func (e *FineTuningEndpoint) CreateFineTuningJob(req *CreateFineTuningJobRequest) (*FineTuningJob, error) {
	var fineTuningJob FineTuningJob
	err := e.do(e, "POST", "jobs", req, nil, &fineTuningJob)
	return &fineTuningJob, err
}

type FineTuningJobs struct {
	Object  string          `json:"object"`
	Data    []FineTuningJob `json:"data"`
	HasMore bool            `json:"has_more"`
}

// Returns a list of paginated fine-tuning job objects.
func (e *FineTuningEndpoint) ListFineTuningJobs(after *string, limit *int) ([]FineTuningJob, error) {
	v := url.Values{}
	if after != nil {
		v.Add("after", *after)
	}
	if limit != nil {
		v.Add("limit", strconv.Itoa(*limit))
	}
	var fineTuningJobs FineTuningJobs
	err := e.do(e, "GET", "jobs", nil, v, &fineTuningJobs)
	// TODO: This needs to move somewhere central
	if err == nil && fineTuningJobs.Object != "list" {
		err = fmt.Errorf("expected 'list' object type, got %s", fineTuningJobs.Object)
	}
	return fineTuningJobs.Data, err
}

// Get info about a fine-tuning job.
// Returns the fine-tuning object with the given ID.
func (e *FineTuningEndpoint) GetFineTuningJob(fineTuningJobId string) (*FineTuningJob, error) {
	var fineTuningJob FineTuningJob
	err := e.do(e, "GET", "jobs/"+fineTuningJobId, nil, nil, &fineTuningJob)
	return &fineTuningJob, err
}

// Immediately cancel a fine-tune job.
// Returns the cancelled fine-tuning object.
func (e *FineTuningEndpoint) CancelFineTuningJob(fineTuningJobId string) (*FineTuningJob, error) {
	var fineTuningJob FineTuningJob
	err := e.do(e, "POST", "jobs/"+fineTuningJobId+"/cancel", nil, nil, &fineTuningJob)
	return &fineTuningJob, err
}

type ListFineTuningEventsRequest struct {
	// Identifier for the last job from the previous pagination request.
	After string `json:"after,omitempty"`

	// Number of fine-tuning jobs to retrieve.
	// Defaults to 20
	Limit int64 `json:"limit,omitempty"`
}

// Get status updates for a fine-tuning job.
// Returns a list of fine-tuning event objects.
func (e *FineTuningEndpoint) ListFineTuningEvents(fineTuningJobId string, req ListFineTuningEventsRequest) ([]FineTuningEvent, error) {
	var fineTuningEvents FineTuningEvents
	err := e.do(e, "GET", "jobs/"+fineTuningJobId+"/events", req, nil, &fineTuningEvents)
	return fineTuningEvents.Data, err
}
