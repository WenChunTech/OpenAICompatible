package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/WenChunTech/OpenAICompatible/src/constant"
)

type RequestBuilder struct {
	*http.Request
}

type Response struct {
	*http.Response
}

func NewRequestBuilder(url, method string) *RequestBuilder {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		slog.Error("Failed to create new request", "error", err)
	}

	return &RequestBuilder{
		Request: req,
	}
}

func (r *RequestBuilder) WithQuery(key, value string) *RequestBuilder {
	query := r.URL.Query()
	query.Add(key, value)
	r.URL.RawQuery = query.Encode()
	return r
}

func (r *RequestBuilder) WithQueries(queries map[string]string) *RequestBuilder {
	query := r.URL.Query()
	for key, value := range queries {
		query.Add(key, value)
	}
	r.URL.RawQuery = query.Encode()
	return r
}

func (r *RequestBuilder) WithHeader(h http.Header) *RequestBuilder {
	r.Header = h
	return r
}

func (r *RequestBuilder) WithHeaders(headers map[string]string) *RequestBuilder {
	for key, value := range headers {
		r.Header.Add(key, value)
	}
	return r
}

func (r *RequestBuilder) WithForm(body map[string][]string) *RequestBuilder {
	r.Header.Set(constant.ContentType, constant.ContentTypeForm)
	r.Body = io.NopCloser(bytes.NewBufferString(url.Values(body).Encode()))
	return r
}

func (r *RequestBuilder) WithJson(body io.Reader) *RequestBuilder {
	r.Header.Set(constant.ContentType, constant.ContentTypeJson)
	r.Body = io.NopCloser(body)
	return r
}

func (r *RequestBuilder) WithFormData(payload io.Reader) *RequestBuilder {
	r.Header.Set(constant.ContentType, constant.ContentTypeFormData)
	r.Body = io.NopCloser(payload)
	return r
}

func (r *RequestBuilder) WithBody(body io.ReadCloser) *RequestBuilder {
	r.Body = body
	return r
}

func (r *RequestBuilder) EnableSSE() *RequestBuilder {
	r.Header.Set(constant.Accept, constant.ContentTypeEventStream)
	r.Header.Set(constant.CacheControl, constant.CacheControlNoCache)
	r.Header.Set(constant.Connection, constant.ConnectionKeepAlive)
	return r
}

func (r *RequestBuilder) Do(ctx context.Context, client *http.Client) (*Response, error) {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(r.Request.WithContext(ctx))
	if err != nil {
		slog.Error("Failed to execute request", "error", err)
		return nil, err
	}
	return &Response{Response: resp}, nil
}

func (r *Response) Json(v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		slog.Error("Failed to decode response body", "error", err)
		return err
	}

	return nil
}

func (r *Response) Text() (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(body), nil
}

func (r *Response) Bytes() ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

// EventStream reads data from the response body as a server-sent event (SSE) stream.
// It returns two channels:
// - A read-only channel for the data chunks ([]byte).
// - A read-only channel for any errors encountered during reading.
// Note: The signature of this method has been changed to include an error channel.
// Calling code may need to be updated to handle the new return values.
func (r *Response) EventStream() (<-chan []byte, <-chan error) {
	dataChan := make(chan []byte)
	errChan := make(chan error, 1) // Buffered channel to avoid blocking

	go func() {
		defer close(dataChan)
		defer close(errChan)
		defer r.Body.Close() // Ensure the response body is closed

		if r.StatusCode != http.StatusOK {
			msg, _ := io.ReadAll(r.Body)
			errMsg := string(msg)
			slog.Error("request failed", "err_msg", errMsg)
			errChan <- fmt.Errorf("status code: %d, err_msg: %s", r.StatusCode, errMsg)
		}

		buf := make([]byte, constant.BufferSize)
		for {
			n, err := r.Body.Read(buf)
			if n > 0 {
				// Create a copy of the slice to avoid data races,
				// as the buffer will be reused.
				data := make([]byte, n)
				copy(data, buf[:n])
				dataChan <- data
			}
			if err != nil {
				if err != io.EOF {
					slog.Error("Failed to read from response body", "error", err)
					errChan <- err
				}
				break // Exit loop on any error, including EOF
			}
		}
	}()

	return dataChan, errChan
}
