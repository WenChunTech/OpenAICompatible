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
	r.Header.Set(constant.ContentType, constant.ContentTypeEventStream)
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

func (r *Response) EventStream() chan []byte {
	ret := make(chan []byte)
	r.Header.Set(constant.ContentType, constant.ContentTypeEventStream)
	r.Header.Set(constant.CacheControl, constant.CacheControlNoCache)
	r.Header.Set(constant.Connection, constant.ConnectionKeepAlive)
	go func() {
		defer close(ret)
		buf := make([]byte, constant.BUFFER_SIZE)
		for {
			n, err := r.Body.Read(buf)
			if err != nil && err != io.EOF {
				slog.Error("Failed to read from response body", "error", err)
			}
			if n == 0 {
				break
			}
			ret <- buf[:n]
		}
	}()

	return ret
}
