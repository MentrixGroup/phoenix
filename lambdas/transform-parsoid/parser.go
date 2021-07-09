package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/wikimedia/phoenix/common"
)

func parseParsoidDocument(document *goquery.Document, snsClient *sns.SNS) (*common.Page, []common.Section, []common.Citations, *common.Citations, *common.Section, error) {
	var err error
	var page *common.Page
	var nodes []common.Section
	var ref *common.Section
	var cts []common.Citations
	var citations *common.Citations

	if page, err = parseParsoidDocumentPage(document); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if nodes, cts, err = parseParsoidDocumentNodes(document, page, snsClient); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if citations, ref, err = parseParsoidDocumentCitation(document, page); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Page, Nodes, Nodes Citations, Page Citations Enhanced, Page Citations Node.
	return page, nodes, cts, citations, ref, nil
}
