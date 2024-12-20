// SPDX-License-Identifier: MPL-2.0

package timeplus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type resource interface {
	resourcePath() string
	resourceID() string
}

type Client struct {
	*http.Client

	baseURL *url.URL
	header  http.Header
}

// optional configurations for the client
type ClientOptions struct {
	BaseURL string
}

func (o *ClientOptions) merge(other ClientOptions) {
	if other.BaseURL != "" {
		o.BaseURL = other.BaseURL
	}
}

func DefaultOptions() ClientOptions {
	return ClientOptions{
		BaseURL: "https://us.timeplus.cloud",
	}
}

func NewClient(workspaceID string, apiKey, username, password string, opts ClientOptions) (*Client, error) {
	ops := DefaultOptions()
	ops.merge(opts)

	baseURL, err := url.Parse(ops.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid BaseURL `%s`: %w", ops.BaseURL, err)
	}
	baseURL = baseURL.JoinPath(workspaceID, "api", "v1beta2")

	return &Client{
		Client:  http.DefaultClient,
		baseURL: baseURL,
		header:  NewHeader(apiKey, username, password),
	}, nil
}

func (c *Client) get(res resource) error {
	req, err := c.newRequest(http.MethodGet, c.baseURL.JoinPath(res.resourcePath(), res.resourceID()).String(), nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}
	return c.do(req, res)
}

func (c *Client) post(res resource) error {
	payload, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("unable to encode request body: %w", err)
	}

	req, err := c.newRequest(http.MethodPost, c.baseURL.JoinPath(res.resourcePath()).String(), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}
	return c.do(req, res)
}

func (c *Client) put(res resource) error {
	payload, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("unable to encode request body: %w", err)
	}

	req, err := c.newRequest(http.MethodPut, c.baseURL.JoinPath(res.resourcePath(), res.resourceID()).String(), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}
	return c.do(req, res)
}

func (c *Client) patch(res resource) error {
	payload, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("unable to encode request body: %w", err)
	}

	req, err := c.newRequest(http.MethodPatch, c.baseURL.JoinPath(res.resourcePath(), res.resourceID()).String(), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}
	return c.do(req, nil)
}

func (c *Client) delete(res resource) error {
	req, err := c.newRequest(http.MethodDelete, c.baseURL.JoinPath(res.resourcePath(), res.resourceID()).String(), nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}
	return c.do(req, nil)
}

func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = c.header
	return req, nil
}

func (c *Client) do(req *http.Request, obj any) error {
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("unable to send request: %w", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("request failed: statusCode=%d body='%s'", resp.StatusCode, string(bodyBytes))
	}

	if obj != nil {
		if err := json.Unmarshal(bodyBytes, obj); err != nil {
			return fmt.Errorf("unable to decode response body %q: %w", string(bodyBytes), err)
		}
	}

	return nil
}
