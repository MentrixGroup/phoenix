package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/wikimedia/phoenix/common"
)

func parseParsoidDocument(document *goquery.Document) (*common.Page, []common.Node, []common.Citations, *common.Citations, *common.Node, error) {
	var err error
	var page *common.Page
	var nodes []common.Node
	var ref *common.Node
	var cts []common.Citations
	var citations *common.Citations

	if page, err = parseParsoidDocumentPage(document); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if nodes, cts, err = parseParsoidDocumentNodes(document, page); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if citations, ref, err = parseParsoidDocumentCitation(document, page); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Page, Nodes, Nodes Citations, Page Citations Enhanced, Page Citations Node.
	return page, nodes, cts, citations, ref, nil
}
