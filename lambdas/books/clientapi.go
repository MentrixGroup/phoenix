package main

import (
	"encoding/json"
	"errors"
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

func (c *Client) GetWikitext(title string) (string, error) {
	var body []byte
	var err error

	page := &WikiPage{}
	query := fmt.Sprintf("/w/api.php?action=parse&format=json&page=%s&prop=wikitext&formatversion=2", title)

	if body, err = c.get(query); err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, page); err != nil {
		return "", err
	}

	return page.Parse.Wikitext, nil
}

func (c *Client) GetCitoidBook(isbn string) (*CitoidBook, error) {
	var body []byte
	var err error
	cbook := &CitoidBooks{}

	query := fmt.Sprintf("/api/rest_v1/data/citation/mediawiki/%s", isbn)
	if body, err = c.get(query); err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := json.Unmarshal(body, &cbook.Books); err != nil {
		return nil, err
	}

	if len(cbook.Books) <= 0 {
		return nil, errors.New("book not found")
	}

	return &cbook.Books[0], nil
}
