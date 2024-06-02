package openai

import (
	"fmt"
)

const VectorStoresEndpointPath = "/vector_stores/"

// VectorStores Endpoint
//
// Vector stores are used to store files for use by the file_search tool.
type VectorStoresEndpoint struct {
	*betaEndpoint
}

// VectorStores Endpoint
func (c *Client) VectorStores() *VectorStoresEndpoint {
	return &VectorStoresEndpoint{newBetaEndpoint(c, VectorStoresEndpointPath)}
}

type FileCounts struct {
	InProgress int64 `json:"in_progress"`
	Completed  int64 `json:"completed"`
	Failed     int64 `json:"failed"`
	Cancelled  int64 `json:"cancelled"`
	Total      int64 `json:"total"`
}

type ExpiresAfter struct {
	Anchor string `json:"anchor"`
	Days   int64  `json:"days"`
}

type VectorStore struct {
	Id           string            `json:"id"`
	Object       string            `json:"object"`
	CreatedAt    int64             `json:"created_at"`
	Name         string            `json:"name"`
	UsageBytes   int64             `json:"usage_bytes"`
	FileCounts   *FileCounts       `json:"file_counts"`
	Status       string            `json:"status"`
	ExpiresAfter ExpiresAfter      `json:"expires_after"`
	ExpiresAt    int64             `json:"expires_at"`
	LastActiveAt int64             `json:"last_active_at"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

type VectorStores struct {
	Object string        `json:"object"`
	Data   []VectorStore `json:"data"`
}

// Returns a list of vector stores.
func (e *VectorStoresEndpoint) ListVectorStores() ([]VectorStore, error) {
	var vectorStores VectorStores
	err := e.do(e, "GET", "", nil, nil, &vectorStores)
	if err == nil && vectorStores.Object != "list" {
		err = fmt.Errorf("expected 'list' object type, got %s", vectorStores.Object)
	}
	return vectorStores.Data, err
}

type CreateVectorStoresRequest struct {
	FileIDs      []string          `json:"file_ids"`
	Name         string            `json:"name"`
	ExpiresAfter *ExpiresAfter     `json:"expires_after,omitempty"`
	MetaData     map[string]string `json:"metadata,omitempty"`
}

// Create a vector store.
func (e *VectorStoresEndpoint) CreateVectorStore(req *CreateVectorStoresRequest) (*VectorStore, error) {
	var vectorStore VectorStore
	err := e.do(e, "POST", "", req, nil, &vectorStore)
	return &vectorStore, err
}

// Retrieves a vector store.
func (e *VectorStoresEndpoint) RetrieveVectorStore(vectorStoreId string) (*VectorStore, error) {
	var vectorStore VectorStore
	err := e.do(e, "GET", vectorStoreId, nil, nil, &vectorStore)
	return &vectorStore, err
}

type ModifyVectorStoresRequest struct {
	Name         string            `json:"name"`
	ExpiresAfter ExpiresAfter      `json:"expires_after"`
	MetaData     map[string]string `json:"metadata,omitempty"`
}

// Modifies a vector store.
func (e *VectorStoresEndpoint) ModifyVectorStore(vectorStoreId string, req *ModifyVectorStoresRequest) (*VectorStore, error) {
	var vectorStore VectorStore
	err := e.do(e, "POST", vectorStoreId, req, nil, &vectorStore)
	return &vectorStore, err
}

type DeletionStatus struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// Deletes a vector store.
func (e *VectorStoresEndpoint) DeleteVectorStore(vectorStoreId string) (*DeletionStatus, error) {
	var status DeletionStatus
	err := e.do(e, "DELETE", vectorStoreId, nil, nil, &status)
	return &status, err
}
