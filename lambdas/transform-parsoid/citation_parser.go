package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/wikimedia/phoenix/common"
)

const slct = "ol.references"

func getRefs(li *goquery.Selection, page string) []string {
	links := make([]string, 0)

	anchors := li.Find("a")

	for i := range anchors.Nodes {
		href := anchors.Eq(i).AttrOr("href", "")
		prefix := fmt.Sprintf("./%s", replaceSpaces(page))

		if ok := strings.HasPrefix(href, prefix); !ok {
			continue
		}

		links = append(links, strings.Replace(href, prefix, "", 1))
	}

	return links
}

// Add citations here when data structure would be ready
func parseParsoidDocumentCitation(document *goquery.Document, page *common.Page) (*common.Citations, *common.Node, error) {
	var err error
	var unsafe string
	var node = &common.Node{}

	citations := &common.Citations{}
	references := document.Find(slct)
	refSec := references.Closest("section").Eq(0)

	for i := range references.Nodes {
		refList := references.Eq(i)
		litems := refList.Find("li")

		for j := range litems.Nodes {
			li := litems.Eq(j)
			isbn := li.Find("bdi").First().Text()

			if unsafe, err = li.Html(); err != nil {
				fmt.Println("error during getting list item HTML content")
				continue
			}

			citation := common.Citation{
				Identifier: li.AttrOr("id", ""),
				Text:       unsafe,
				References: getRefs(li, page.Name),
			}

			if isbn != "" && validate(isbn) {
				citation.Source = getSourceId(isbn)
			}

			citations.Citations = append(citations.Citations, citation)
		}
	}

	node.Name = getSectionName(refSec)
	node.DateModified = page.DateModified

	if unsafe, err = refSec.Html(); err != nil {
		return citations, nil, err
	}

	citations.IsPartOf = []string{fmt.Sprintf("pages/%s/%s_citations_enhanced", replaceSpaces(page.Name), replaceSpaces(page.Name))}
	node.ID = fmt.Sprintf("pages/%s/%s_citations", replaceSpaces(page.Name), replaceSpaces(page.Name))
	node.Unsafe = unsafe

	return citations, node, nil
}
