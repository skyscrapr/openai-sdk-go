package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	apiURL = "https://api.openai.com"
)

// Client - OpenAI client.
type Client struct {
	authToken string

	BaseURL        *url.URL
	OrganizationID string
	HTTPClient     *http.Client
	UserAgent      string
}

// NewClient creates new OpenAI client.
func NewClient(authToken string) *Client {
	c := &Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		authToken:  authToken,
		UserAgent:  "skyscrapr/openai-sdk-go",
	}
	c.BaseURL, _ = url.Parse(apiURL)
	return c
}

func (c *Client) do(e endpointI, method string, path string, body interface{}, result interface{}) error {
	u := e.buildURL(path)
	req, err := e.newRequest(method, u, body)
	if err != nil {
		return err
	}
	return e.doRequest(req, result)
}

func (c *Client) newRequest(method string, u *url.URL, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	// req.Header.Set("User-Agent", c.UserAgent)

	if len(c.OrganizationID) > 0 {
		req.Header.Set("OpenAI-Organization", c.OrganizationID)
	}

	// if body != nil {
	// 	req.Header.Set("Content-Type", "application/json")
	// }

	return req, nil
}

func (c *Client) doRequest(req *http.Request, v any) error {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return c.handleErrorResp(res)
	}

	return decodeResponse(res.Body, v)
}

func decodeResponse(body io.Reader, v any) error {
	if v == nil {
		return nil
	}
	err := json.NewDecoder(body).Decode(v)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (c *Client) handleErrorResp(resp *http.Response) error {
	var errRes ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errRes)
	if err != nil || errRes.Error == nil {
		reqErr := &RequestError{
			HTTPStatusCode: resp.StatusCode,
			Err:            err,
		}
		if errRes.Error != nil {
			reqErr.Err = errRes.Error
		}
		return reqErr
	}

	errRes.Error.HTTPStatusCode = resp.StatusCode
	return errRes.Error
}
