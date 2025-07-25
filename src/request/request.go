package request

import (
	"bytes"
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

func (r *RequestBuilder) AddQuery(key, value string) *RequestBuilder {
	query := r.URL.Query()
	query.Add(key, value)
	r.URL.RawQuery = query.Encode()
	return r
}

func (r *RequestBuilder) AddQueries(queries map[string]string) *RequestBuilder {
	query := r.URL.Query()
	for key, value := range queries {
		query.Add(key, value)
	}
	r.URL.RawQuery = query.Encode()
	return r
}

func (r *RequestBuilder) AddHeader(key, value string) *RequestBuilder {
	r.Header.Add(key, value)
	return r
}

func (r *RequestBuilder) AddHeaders(headers map[string]string) *RequestBuilder {
	for key, value := range headers {
		r.Header.Add(key, value)
	}
	return r
}

func (r *RequestBuilder) SetForm(body map[string][]string) *RequestBuilder {
	r.Header.Set(constant.ContentType, constant.ContentTypeForm)
	r.Body = io.NopCloser(bytes.NewBufferString(url.Values(body).Encode()))
	return r
}

func (r *RequestBuilder) SetJson(body io.Reader) *RequestBuilder {
	r.Header.Set(constant.ContentType, constant.ContentTypeJson)
	r.Body = io.NopCloser(body)
	return r
}

func (r *RequestBuilder) SetFormData(payload io.Reader) *RequestBuilder {
	r.Header.Set(constant.ContentType, constant.ContentTypeFormData)
	r.Body = io.NopCloser(payload)
	return r
}

func (r *RequestBuilder) EnableSSE() *RequestBuilder {
	r.Header.Set(constant.Accept, constant.ContentTypeEventStream)
	r.Header.Set(constant.CacheControl, constant.CacheControlNoCache)
	r.Header.Set(constant.Connection, constant.ConnectionKeepAlive)
	return r
}

func (r *RequestBuilder) Do(client *http.Client) (*Response, error) {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(r.Request)
	if err != nil {
		slog.Error("Failed to execute request", "error", err)
		return nil, err
	}
	return &Response{Response: resp}, nil
}

func (r *Response) Json(v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	defer r.Body.Close()
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

func (r *Response) Text() (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	defer r.Body.Close()
	return string(body), nil
}

func (r *Response) Bytes() ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	defer r.Body.Close()
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

		buf := make([]byte, constant.BUFFER_SIZE)
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
