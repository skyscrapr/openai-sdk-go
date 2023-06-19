// Package sseclient is library for consuming SSE streams.
//
// Key features:
//
// Synchronous execution. Reconnecting, event parsing and processing is executed
// in single go-routine that started the stream. This gives freedom to use any
// concurrency and synchronization model.
//
// Go context aware. SSE streams can be optionally given a context on start.
// This gives flexibility to support different stream stopping mechanisms.
package openai

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

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

// SSEClient is used to connect to SSE stream and receive events. It handles HTTP
// request creation and reconnects automatically.
//
// There appears to be an issue with the OpenAI SSE event server.
// No LAST-EVENT-ID is supplied. Therefore resuming after reconnection is difficult.
// Reconnection will generally start the stream from the beginning of the event stream.
// Therefore for simplicity we use a very simple mechanism to count processed events and
// skip the pre-seen events on reconnection.
type SSEClient struct {
	URL                 string
	LastEventID         string
	Retry               time.Duration
	HTTPClient          *http.Client
	Headers             http.Header
	processedEventCount int

	// VerboseStatusCodes specifies whether connect should return all
	// status codes as errors if they're not StatusOK (200).
	VerboseStatusCodes bool
}

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

// New creates SSE stream client object. It will use given url and
// last event ID values and a 2 second retry timeout.
// It will use custom http client that skips verification for tls process.
// This method only creates Client struct and does not start connecting to the
// SSE endpoint.
func NewSSEClient(url, lastEventID string) *SSEClient {
	return &SSEClient{
		URL:         url,
		LastEventID: lastEventID,
		Retry:       2 * time.Second,
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		Headers:             make(http.Header),
		processedEventCount: 0,
	}
}

// StreamMessage stores single SSE event or error.
type StreamMessage struct {
	Event *SSEEvent
	Err   error
}

type SSEEvent struct {
	Data string
}

// Stream is non-blocking SSE stream consumption mode where events are passed
// through a channel. Stream can be stopped by cancelling context.
//
// Parameter buf controls returned stream channel buffer size. Buffer size of 0
// is a good default.
func (c *SSEClient) Stream(ctx context.Context, buf int) <-chan StreamMessage {
	ch := make(chan StreamMessage, buf)
	errorFn := func(err error) error {
		select {
		case ch <- StreamMessage{Err: err}:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	eventFn := func(e *SSEEvent) error {
		select {
		case ch <- StreamMessage{Event: e}:
		case <-ctx.Done():
		}
		return nil
	}

	go func() {
		defer close(ch)
		c.Start(ctx, eventFn, errorFn)
	}()

	return ch
}

// Start connects to the SSE stream. This function will block until SSE stream
// is stopped. Stopping SSE stream is possible by cancelling given stream
// context or by returning some error from the error handler callback. Error
// returned by the error handler is passed back to the caller of this function.
func (c *SSEClient) Start(ctx context.Context, eventFn EventHandler, errorFn EventErrorHandler) error {
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
		err := c.connect(ctx, eventFn)
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

// connect performs single connection to SSE endpoint.
func (c *SSEClient) connect(ctx context.Context, eventFn EventHandler) error {
	eventCount := 0
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.URL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")

	for h, vs := range c.Headers {
		for _, v := range vs {
			req.Header.Add(h, v)
		}
	}

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
				c.processedEventCount++
				if err := eventFn(event); err != nil {
					return err
				}
			}
		}
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		if c.VerboseStatusCodes {
			return fmt.Errorf("bad response status code %d", resp.StatusCode)
		}
		// reconnect without logging an error.
		return errStreamConn
	default:
		// trigger a reconnect and output an error.
		return fmt.Errorf("bad response status code %d", resp.StatusCode)
	}
}

// parseEvent reads a single Event fromthe event stream.
func (c *SSEClient) parseEvent(r *bufio.Reader) (*SSEEvent, error) {

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
