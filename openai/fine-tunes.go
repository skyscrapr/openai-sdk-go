package openai

import (
	"context"
	"fmt"
	"time"
)

const FineTunesEndpointPath = "/fine-tunes/"

// FineTunes Endpoint
//
// Manage fine-tuning jobs to tailor a model to your specific training data.
// Related guide: [Fine-tune models]: https://platform.openai.com/docs/guides/fine-tuning
type FineTunesEndpoint struct {
	*endpoint
}

// Files Endpoint
func (c *Client) FineTunes() *FineTunesEndpoint {
	return &FineTunesEndpoint{newEndpoint(c, FineTunesEndpointPath)}
}

type FineTune struct {
	Id        string `json:"id"`
	Object    string `json:"object"`
	Model     string `json:"model"`
	CreatedAt int64  `json:"created_at"`
	Events    []struct {
		Object    string `json:"object"`
		CreatedAt int64  `json:"created_at"`
		Level     string `json:"level"`
		Message   string `json:"message"`
	} `json:"events"`
	FineTunedModel string `json:"fine_tuned_model"`
	Hyperparams    struct {
		BatchSize              int64   `json:"batch_size"`
		LearningRateMultiplier float64 `json:"learning_rate_multiplier"`
		NEpochs                int64   `json:"n_epochs"`
		PromptLossWeight       float64 `json:"prompt_loss_weight"`
	} `json:"hyperparams"`
	OrganizationId  string `json:"organization_id"`
	ResultFiles     []File `json:"result_files"`
	Status          string `json:"status"`
	ValidationFiles []File `json:"validation_files"`
	TrainingFiles   []File `json:"training_files"`
	UpdatedAt       int64  `json:"updated_at"`
}

type FineTuneEvents struct {
	Object string          `json:"object"`
	Data   []FineTuneEvent `json:"data"`
}

type FineTuneEvent struct {
	Object    string `json:"object"`
	CreatedAt int    `json:"created_at"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type CreateFineTunesRequest struct {
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
	// Defaults to curie
	// The name of the base model to fine-tune. You can select one of "ada", "babbage", "curie", "davinci", or a fine-tuned model created after 2022-04-21. To learn more about these models, see the Models documentation.
	Model string `json:"model,omitempty"`
	// Defaults to 4
	// The number of epochs to train the model for. An epoch refers to one full cycle through the training dataset.
	NEpochs int64 `json:"n_epochs,omitempty"`
	// Defaults to null
	// The batch size to use for training. The batch size is the number of training examples used to train a single forward and backward pass.
	// By default, the batch size will be dynamically configured to be ~0.2% of the number of examples in the training set, capped at 256 - in general, we've found that larger batch sizes tend to work better for larger datasets.
	BatchSize int64 `json:"batch_size,omitempty"`
	// Defaults to null
	// The learning rate multiplier to use for training. The fine-tuning learning rate is the original learning rate used for pretraining multiplied by this value.
	// By default, the learning rate multiplier is the 0.05, 0.1, or 0.2 depending on final batch_size (larger learning rates tend to perform better with larger batch sizes). We recommend experimenting with values in the range 0.02 to 0.2 to see what produces the best results.
	LearningRateMultiplier int64 `json:"learning_rate_multiplier,omitempty"`
	// Defaults to 0.01
	// The weight to use for loss on the prompt tokens. This controls how much the model tries to learn to generate the prompt (as compared to the completion which always has a weight of 1.0), and can add a stabilizing effect to training when completions are short.
	// If prompts are extremely long (relative to completions), it may make sense to reduce this weight so as to avoid over-prioritizing learning the prompt.
	PromptLossWeight int64 `json:"prompt_loss_weight,omitempty"`
	// 	Defaults to false
	// If set, we calculate classification-specific metrics such as accuracy and F-1 score using the validation set at the end of every epoch. These metrics can be viewed in the results file.
	// In order to compute classification metrics, you must provide a validation_file. Additionally, you must specify classification_n_classes for multiclass classification or classification_positive_class for binary classification.
	ComputeClassificationMetrics bool `json:"compute_classification_metrics,omitempty"`
	// Defaults to null
	// The number of classes in a classification task.
	// This parameter is required for multiclass classification.
	ClassificationNClasses int64 `json:"classification_n_classes,omitempty"`
	// The positive class in binary classification.
	// This parameter is needed to generate precision, recall, and F1 metrics when doing binary classification.
	ClassificationPositiveClass string `json:"classification_positive_class,omitempty"`
	// If this is provided, we calculate F-beta scores at the specified beta values. The F-beta score is a generalization of F-1 score. This is only used for binary classification.
	// With a beta of 1 (i.e. the F-1 score), precision and recall are given the same weight. A larger beta score puts more weight on recall and less on precision. A smaller beta score puts more weight on precision and less on recall.
	ClassificationBetas []string `json:"classification_betas,omitempty"`
	// 	A string of up to 40 characters that will be added to your fine-tuned model name.
	// For example, a suffix of "custom-model-name" would produce a model name like ada:ft-your-org:custom-model-name-2022-02-15-04-21-04.
	Suffix []string `json:"suffix,omitempty"`
}

// Creates a job that fine-tunes a specified model from a given dataset.
// Response includes details of the enqueued job including job status and the name of the fine-tuned models once complete.
// Learn more about Fine-tuning
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/fine-tunes
func (e *FineTunesEndpoint) CreateFineTune(req *CreateFineTunesRequest) (*FineTune, error) {
	var fineTune FineTune
	err := e.do(e, "POST", "", req, &fineTune)
	return &fineTune, err
}

type FineTunes struct {
	Object string     `json:"object"`
	Data   []FineTune `json:"data"`
}

// List your organization's fine-tuning jobs
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/fine-tunes
func (e *FineTunesEndpoint) ListFineTunes() ([]FineTune, error) {
	var fineTunes FineTunes
	err := e.do(e, "GET", "", nil, &fineTunes)
	// TODO: This needs to move somewhere central
	if err == nil && fineTunes.Object != "list" {
		err = fmt.Errorf("expected 'list' object type, got %s", fineTunes.Object)
	}
	return fineTunes.Data, err
}

// Gets info about the fine-tune job. [Learn more about Fine-tuning]: https://platform.openai.com/docs/guides/fine-tuning
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/fine-tunes
func (e *FineTunesEndpoint) GetFineTune(fineTuneId string) (*FineTune, error) {
	var fineTune FineTune
	err := e.do(e, "GET", fineTuneId, nil, &fineTune)
	return &fineTune, err
}

// Immediately cancel a fine-tune job.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/fine-tunes
func (e *FineTunesEndpoint) CancelFineTune(fineTuneId string) (*FineTune, error) {
	var fineTune FineTune
	err := e.do(e, "POST", fineTuneId+"/cancel", nil, &fineTune)
	return &fineTune, err
}

// Get fine-grained status updates for a fine-tune job.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/fine-tunes
// TODO: Need to add support for streams
func (e *FineTunesEndpoint) ListFineTuneEvents(fineTuneId string) ([]FineTuneEvent, error) {
	var fineTuneEvents FineTuneEvents
	err := e.do(e, "POST", fineTuneId+"/events", nil, fineTuneEvents)
	return fineTuneEvents.Data, err
}

// Get streamed status updates for a fine-tune job.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/fine-tunes
func (e *FineTunesEndpoint) SubscribeFineTuneEvents(fineTuneId string, eventHandler EventHandler, errorHandler EventErrorHandler) error {
	u, err := e.buildURL(fineTuneId + "/events?stream=true")
	if err != nil {
		return err
	}
	req, err := e.newRequest("GET", u, nil)
	if err != nil {
		return err
	}
	// req.Header.Set("Connection", "keep-alive")

	c := NewSSEClient(u.String(), "")
	c.HTTPClient.Timeout = 0
	c.Headers = req.Header
	ctx, _ := context.WithTimeout(context.Background(), time.Hour)

	return c.Start(ctx, eventHandler, errorHandler)
}

// Delete a fine-tuned model. You must have the Owner role in your organization.
// This is actually implemented in the models endpoint
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/fine-tunes/delete-model
func (e *FineTunesEndpoint) DeleteFineTuneModel(modelId string) (bool, error) {
	return e.Client.Models().deleteFineTuneModel(modelId)
}
