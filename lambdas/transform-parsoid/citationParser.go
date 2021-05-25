package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/wikimedia/phoenix/common"
)

const slct = "ol.references"

// Add citations here when data structure would be ready
func parseParsoidDocumentCitation(document *goquery.Document, page *common.Page) (common.Node, error) {
	var err error
	var unsafe string

	node := common.Node{
		Source:       page.Source,
		DateModified: page.DateModified,
		Name:         "References",
	}
	citation := document.Find(slct).First().Closest("section")

	if unsafe, err = citation.Html(); err != nil {
		return common.Node{}, err
	}

	node.Unsafe = unsafe

	fmt.Println(node)

	return node, nil
}
