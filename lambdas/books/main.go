// Package lambda to parse books
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/wikimedia/phoenix/common"
)

type CitoidBooks struct {
	Books []CitoidBook
}

type CitoidBook struct {
	Itemtype    string     `json:"itemType"`
	Date        string     `json:"date"`
	Publisher   string     `json:"publisher"`
	Title       string     `json:"title"`
	Isbn        []string   `json:"ISBN"`
	Place       string     `json:"place"`
	Numpages    string     `json:"numPages"`
	Oclc        string     `json:"oclc"`
	URL         string     `json:"url"`
	Contributor [][]string `json:"contributor"`
	Author      [][]string `json:"author"`
	Accessdate  string     `json:"accessDate"`
	Source      []string   `json:"source"`
}

type Book struct {
	Date      string   `json:"date"`
	Publisher string   `json:"publisher"`
	Title     string   `json:"title"`
	Isbn      string   `json:"isbn"`
	URL       string   `json:"url"`
	Author    []string `json:"author"`
	Editor    string   `json:"editor"`
	Numpages  string   `json:"numberOfPages"`
}

type WikiPage struct {
	Parse struct {
		Title    string `json:"title"`
		Pageid   int    `json:"pageid"`
		Wikitext string `json:"wikitext"`
	} `json:"parse"`
}

const (
	left    = "{{cite book"
	right   = "}}"
	baseURL = "https://simple.wikipedia.org"
)

var logger *common.Logger

func regx(p string) string {
	return fmt.Sprintf(`(?s)%s[ \t]*=[ \t]*(.*?)\|`, regexp.QuoteMeta(p))
}

func snmatch(s string, matcher string) []string {
	a := make([]string, 0)
	r := regexp.MustCompile(fmt.Sprintf(`(?s)%s[0-9][ \t]*=[ \t]*(.*?)\|`, matcher))
	matches := r.FindAllStringSubmatch(s, -1)

	if len(matches) > 0 {
		for _, v := range matches {
			a = append(a, strings.TrimSpace(v[1]))
		}

		return a
	}

	return a
}

func smatch(s string, matcher string) string {
	r := regexp.MustCompile(regx(matcher))
	matches := r.FindStringSubmatch(s)

	if len(matches) > 0 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}

func saveBook(b *Book, title string) {
	f, err := os.Create(fmt.Sprintf("%s.json", title))
	if err != nil {
		log.Panic(err)
	}

	defer f.Close()

	if err := json.NewEncoder(f).Encode(b); err != nil {
		log.Panic(err)
	}
}

func Handler(t string, cl *Client) {
	page := &WikiPage{}
	books := make([]Book, 0)
	query := fmt.Sprintf("/w/api.php?action=parse&format=json&page=%s&prop=wikitext&formatversion=2", t)
	// Perform the HTTP request
	var body []byte
	var err error

	if body, err = cl.get(query); err != nil {
		logger.Error("error making HTTP request: %w", err)
	}

	if err := json.Unmarshal(body, page); err != nil {
		logger.Error("error unmarshaling data: %w", err)
	}

	wikitext := page.Parse.Wikitext
	rx := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(left) + `(.*?)` + regexp.QuoteMeta(right))

	matches := rx.FindAllStringSubmatch(wikitext, -1)

	for _, v := range matches {
		cbook := &CitoidBooks{}
		line := fmt.Sprintf("%s%s", strings.TrimSpace(v[1]), "|")

		editor := fmt.Sprintf("%s %s", smatch(line, "first"), smatch(line, "last"))

		book := &Book{
			Date:      smatch(line, "date"),
			Publisher: smatch(line, "publisher"),
			Title:     smatch(line, "title"),
			Isbn:      smatch(line, "isbn"),
			URL:       smatch(line, "url"),
			Author:    snmatch(line, "author"),
			Editor:    editor,
		}

		books = append(books, *book)

		if body, err = cl.get(fmt.Sprintf("/api/rest_v1/data/citation/mediawiki/%s", book.Isbn)); err != nil {
			logger.Error("Error making HTTP request: %w", err)
		}

		if err := json.Unmarshal(body, &cbook.Books); err != nil {
			logger.Error("Error unmarshaling data: %w", err)
		}

		if len(cbook.Books) > 0 {
			book.Numpages = cbook.Books[0].Numpages
		}

		saveBook(book, fmt.Sprintf("%s_%s", t, book.Title))
	}

}

func init() {
	// Determine logging level
	var level string = "ERROR"
	if v, ok := os.LookupEnv("LOG_LEVEL"); ok {
		level = v
	}

	// Initialize the logger
	logger = common.NewLogger(level)
}

func main() {
	cl := NewClient(baseURL)

	Handler("Albert_Einstein", cl)
}
