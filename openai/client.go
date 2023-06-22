package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiURL = "https://api.openai.com"
)

// Client - OpenAI client.
type Client struct {
	authToken string

	BaseURL             *url.URL
	OrganizationID      string
	HTTPClient          *http.Client
	UserAgent           string
	Retry               time.Duration
	processedEventCount int
}

// NewClient creates new OpenAI client.
func NewClient(authToken string) *Client {
	c := &Client{
		HTTPClient:          &http.Client{Timeout: 30 * time.Second},
		authToken:           authToken,
		UserAgent:           "skyscrapr/openai-sdk-go",
		Retry:               2 * time.Second,
		processedEventCount: 0,
	}
	c.BaseURL, _ = url.Parse(apiURL)
	return c
}

func (c *Client) do(e endpointI, method string, path string, body interface{}, result interface{}) error {
	u, err := e.buildURL(path)
	if err != nil {
		return err
	}
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
	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	c.setCommonHeaders(req)

	return req, nil
}

func (c *Client) setCommonHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	if len(c.OrganizationID) > 0 {
		req.Header.Set("OpenAI-Organization", c.OrganizationID)
	}
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

// EventErrorHandler is a callback that gets called every time SSE stream encounters
// an error including errors returned by EventHandler function. Network
// connection errors and response codes 500, 502, 503, 504 are not treated as
// errors.
//
// If error handler returns nil, error will be treated as handled and stream
// will continue to be processed (with automatic reconnect).
//
// If error handler returns error it is treated as fatal and stream processing
// loop exits returning received error up the stack.
//
// This handler can be used to implement complex error handling scenarios. For
// simple cases ReconnectOnError or StopOnError are provided by this library.
//
// Users of this package have to provide this function implementation.
type EventErrorHandler func(error) error

// EventHandler is a callback that gets called every time event on the SSE
// stream is received. Error returned from handler function will be passed to
// the error handler.
//
// Users of this package have to provide this function implementation.
type EventHandler func(e *SSEEvent) error

// List of commonly used error handler function implementations.
var (
	ReconnectOnError EventErrorHandler = func(error) error { return nil }
	StopOnError      EventErrorHandler = func(err error) error { return err }
)

var (
	// MalformedEvent error is returned if stream ended with incomplete event.
	ErrorMalformedEvent = errors.New("incomplete event at the end of the stream")

	// errStreamConn error is returned when client is unable to
	// connect to the stream. This error is only used to reconnect to
	// the stream without outputing connection errors to the client.
	errStreamConn = errors.New("cannot connect to the stream")
)

type SSEEvent struct {
	Data string
}

// stream connects to the SSE stream. This function will block until SSE stream
// is stopped. Stopping SSE stream is possible by cancelling given stream
// context or by returning some error from the error handler callback. Error
// returned by the error handler is passed back to the caller of this function.
func (c *Client) stream(ctx context.Context, u string, eventFn EventHandler, errorFn EventErrorHandler) error {
	c.processedEventCount = 0
	lastTimeout := c.Retry / 32
	tm := time.NewTimer(0)
	stop := func() {
		tm.Stop()

		select {
		case <-tm.C:
		default:
		}
	}
	defer stop()

	for {
		err := c.connectSSE(ctx, u, eventFn)
		switch err {
		case io.EOF, io.ErrUnexpectedEOF:
			// ok, we can reconnect right away
			lastTimeout = c.Retry / 32
		case nil, ctx.Err():
			// context cancellation exits silently
			return nil
		default:
			if !errors.Is(err, errStreamConn) {
				if cerr := errorFn(err); cerr != nil {
					// error handler instructs to stop
					// the sse stream
					return cerr
				}
			}

			stop()
			tm.Reset(lastTimeout)

			select {
			case <-tm.C:
			case <-ctx.Done():
				// context cancellation exits silently
				return nil
			}

			if lastTimeout < c.Retry {
				lastTimeout = lastTimeout * 2
			}
		}
	}
}

// connectSSE performs single connection to SSE endpoint.
func (c *Client) connectSSE(ctx context.Context, u string, eventFn EventHandler) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	c.setCommonHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// silently ignore connection errors and reconnect.
		return errStreamConn
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// we do not support BOM in sse streams, or \r line separators.
		r := bufio.NewReader(resp.Body)

		eventCount := 0
		for {
			event, err := c.parseEvent(r)
			if err != nil {
				return err
			}
			if event == nil {
				return nil
			}

			eventCount++
			if eventCount > c.processedEventCount {
				if err := eventFn(event); err != nil {
					return err
				}
				c.processedEventCount++
			}
		}
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		// reconnect without logging an error.
		return errStreamConn
	default:
		// trigger a reconnect and output an error.
		return fmt.Errorf("bad response status code %d", resp.StatusCode)
	}
}

// parseEvent reads a single Event fromthe event stream.
func (c *Client) parseEvent(r *bufio.Reader) (*SSEEvent, error) {

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			// EOF is treated as silent reconnect. If this is
			// malformed event report an error.
			if err == io.EOF && len(line) != 0 {
				err = ErrorMalformedEvent
			}
			return nil, err
		}
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "data:") {
			line = strings.TrimPrefix(line, "data:")
			line = strings.TrimSpace(line)
			if line == "[DONE]" {
				return nil, nil
			}
			event := SSEEvent{
				Data: line,
			}
			return &event, nil
		}
		continue
	}
}
