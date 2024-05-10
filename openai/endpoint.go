package openai

import (
	"net/http"
	"net/url"
	"path"
)

const (
	apiPath = "v1"
)

type endpointI interface {
	buildURL(endpoint string) (*url.URL, error)
	newRequest(method string, u *url.URL, body interface{}) (*http.Request, error)
	doRequest(req *http.Request, v any) error
}

type endpoint struct {
	*Client
	EndpointPath string
}

type betaEndpoint struct {
	endpoint
}

func newEndpoint(c *Client, endpointPath string) *endpoint {
	e := &endpoint{
		Client:       c,
		EndpointPath: endpointPath,
	}
	return e
}

func newBetaEndpoint(c *Client, endpointPath string) *betaEndpoint {
	e := &betaEndpoint{
		endpoint: *newEndpoint(c, endpointPath),
	}
	return e
}

func (e *endpoint) buildURL(endpointPath string) (*url.URL, error) {
	u, err := url.Parse(endpointPath)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(e.EndpointPath, u.Path)
	u.Path = path.Join(apiPath, u.Path)
	u.Path = path.Join(e.BaseURL.Path, u.Path)
	return e.BaseURL.ResolveReference(u), err
}

func (e *endpoint) doRequest(req *http.Request, v any) error {
	return e.Client.doRequest(req, v)
}

func (e *endpoint) newRequest(method string, u *url.URL, body interface{}) (*http.Request, error) {
	return e.Client.newRequest(method, u, body)
}

func (e *betaEndpoint) newRequest(method string, u *url.URL, body interface{}) (*http.Request, error) {
	req, err := e.Client.newRequest(method, u, body)
	req.Header.Set("OpenAI-Beta", "assistants=v2")
	return req, err
}
