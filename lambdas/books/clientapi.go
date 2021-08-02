package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/AlisterIgnatius/phoenix/common"
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

type GoogleBook struct {
	Volumeinfo struct {
		Title         string   `json:"title"`
		Subtitle      string   `json:"subtitle"`
		Authors       []string `json:"authors"`
		Publisher     string   `json:"publisher"`
		Publisheddate string   `json:"publishedDate"`
		Imagelinks    struct {
			Smallthumbnail string `json:"smallThumbnail"`
			Thumbnail      string `json:"thumbnail"`
		} `json:"imageLinks"`
	} `json:"volumeInfo"`
}

type GoogleBooksResponse struct {
	Kind       string       `json:"kind"`
	Totalitems int          `json:"totalItems"`
	Items      []GoogleBook `json:"items"`
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

func (c *Client) GetBook(isbn string) (*common.Book, error) {
	var body []byte
	var err error

	response := &GoogleBooksResponse{}
	query := fmt.Sprintf("/books/v1/volumes?q=isbn:%s", isbn)

	if body, err = c.get(query); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	if len(response.Items) <= 0 {
		return nil, fmt.Errorf("error book not found: %w", err)
	}

	gbook := response.Items[0].Volumeinfo

	book := &common.Book{
		Isbn:          isbn,
		Name:          fmt.Sprintf("%s: %s", gbook.Title, gbook.Subtitle),
		Author:        gbook.Authors,
		Publisher:     gbook.Publisher,
		Datepublished: gbook.Publisheddate,
		Thumbnailurl:  gbook.Imagelinks.Thumbnail,
	}

	return book, nil
}
