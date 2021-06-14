package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// NewClient creates new API client
func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: new(http.Client),
		baseURL:    baseURL,
	}
}

// Client API
type Client struct {
	httpClient *http.Client
	baseURL    string
}

func (c *Client) get(query string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.baseURL, query), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET error: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d (expected %d)", res.StatusCode, http.StatusOK)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading HTTP response body: %w", err)
	}

	return body, nil
}
