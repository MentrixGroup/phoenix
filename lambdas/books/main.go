// Package lambda to parse books
package main

import (
	"encoding/json"
	"fmt"
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

func snmatch(str string, matcher string) []string {
	match := make([]string, 0)
	rgx := regexp.MustCompile(fmt.Sprintf(`(?s)%s[0-9][ \t]*=[ \t]*(.*?)\|`, matcher))
	matches := rgx.FindAllStringSubmatch(str, -1)

	if len(matches) > 0 {
		for _, v := range matches {
			match = append(match, strings.TrimSpace(v[1]))
		}

		return match
	}

	return match
}

func smatch(s string, matcher string) string {
	r := regexp.MustCompile(fmt.Sprintf(`(?s)%s[ \t]*=[ \t]*(.*?)\|`, regexp.QuoteMeta(matcher)))
	matches := r.FindStringSubmatch(s)

	if len(matches) > 0 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}

func Handler(title string, cl *Client) {
	books := make([]Book, 0)

	var cbook *CitoidBook
	var wikitext string
	var err error

	if wikitext, err = cl.GetWikitext(title); err != nil {
		logger.Error("error making HTTP request: %w", err)
		return
	}

	rx := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(left) + `(.*?)` + regexp.QuoteMeta(right))

	matches := rx.FindAllStringSubmatch(wikitext, -1)

	for _, v := range matches {
		line := fmt.Sprintf("%s%s", strings.TrimSpace(v[1]), "|")

		book := &Book{
			Date:      smatch(line, "date"),
			Publisher: smatch(line, "publisher"),
			Title:     smatch(line, "title"),
			Isbn:      smatch(line, "isbn"),
			URL:       smatch(line, "url"),
			Author:    snmatch(line, "author"),
			Editor:    fmt.Sprintf("%s %s", smatch(line, "first"), smatch(line, "last")),
		}

		books = append(books, *book)

		if cbook, err = cl.GetCitoidBook(book.Isbn); err != nil {
			logger.Error("Error getting citoid books: %w", err)
		}

		if cbook != nil {
			book.Numpages = cbook.Numpages
		}

		f, err := os.Create(fmt.Sprintf("%s.json", fmt.Sprintf("%s_%s", title, book.Title)))
		if err != nil {
			logger.Error("Error creating JSON file: %w", err)
		}

		defer f.Close()

		if err := json.NewEncoder(f).Encode(book); err != nil {
			logger.Error("Error encoding JSON file: %w", err)
		}
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
