package main

import (
	"fmt"
	"strings"

	"github.com/AlisterIgnatius/phoenix/common"
	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go/service/sns"
)

var (
	ignoredNodes = map[string]bool{
		"References": true,
	}
)

func getSectionName(section *goquery.Selection) string {
	return section.Find("h2").First().Text()
}

func getCitationIds(section *goquery.Selection, page string) []string {
	ids := make([]string, 0)
	links := section.Find("sup.reference")

	for i := range links.Nodes {
		refNode := links.Eq(i)
		ids = append(ids, strings.Replace(refNode.Find("a").Eq(0).AttrOr("href", ""), fmt.Sprintf("./%s#", replaceSpaces(page)), "", 1))
	}

	return ids
}

func getCitation(document *goquery.Document, cId, page string, section string, snsClient *sns.SNS) (common.Citation, error) {
	var err error
	var unsafe string

	citation := common.Citation{}
	slct := document.Find(fmt.Sprintf("#%s", cId))

	if len(slct.Nodes) <= 0 {
		return citation, fmt.Errorf("No element was found with id=%s found", cId)
	}

	li := slct.Eq(0)
	isbn := li.Find("bdi").First().Text()

	if unsafe, err = li.Html(); err != nil {
		return citation, fmt.Errorf("error during getting list item HTML content")
	}

	citation.Identifier = li.AttrOr("id", "")
	citation.References = getRefs(li, page)
	citation.Text = unsafe

	if isbn != "" && validate(isbn) {
		citation.Source = getSourceId(isbn)
		output, err := sourceParseEvent(snsClient, &common.SourseParseEvent{ID: isbn, Page: page, Section: section})

		if err != nil {
			return citation, fmt.Errorf("error sendig out SNS event: %s", err)
		}

		log.Debug("Successfully published SNS message: %s", *output.MessageId)
	}

	return citation, nil
}

func parseParsoidDocumentNodes(document *goquery.Document, page *common.Page, snsClient *sns.SNS) ([]common.Section, []common.Citations, error) {
	var err error
	var modified = page.DateModified
	var nameCounts = make(map[string]int)
	var nodes = make([]common.Section, 0)
	var sections = document.Find("html>body>section[data-mw-section-id]")
	var cits = make([]common.Citations, 0)

	for i := range sections.Nodes {
		var node = common.Section{}
		var ct = common.Citations{}
		var section = sections.Eq(i)
		var unsafe string
		var nodeCites = make([]common.Citation, 0)

		// If this is the first section and the name is a zero length string, then we assign it
		// a constant to simplify lookups
		if i == 0 {
			if name := getSectionName(section); name == "" {
				node.Name = leadSectionName
			} else {
				node.Name = name
			}
		} else {
			node.Name = getSectionName(section)
		}

		// Since it is possible for a document to have more than one section with the same heading text, keep
		// track of the number of times we've assigned a name, and de-duplicate if necessary.
		nameCounts[strings.ToLower(node.Name)]++

		if nameCounts[strings.ToLower(node.Name)] > 1 {
			node.Name = fmt.Sprintf("%s_%d", node.Name, nameCounts[strings.ToLower(node.Name)])
		}

		for _, id := range getCitationIds(section, page.Name) {
			ct, err := getCitation(document, id, page.Name, node.Name, snsClient)

			if err != nil {
				log.Error("problem with creating citation: %s", err)
			}
			nodeCites = append(nodeCites, ct)
		}

		ct.Citations = nodeCites

		node.DateModified = modified
		node.Version = page.Version
		unqn := fmt.Sprintf("%s_%s", replaceSpaces(page.Name), replaceSpaces(node.Name))
		node.HasPart = []common.Entity{common.Entity{Identifier: fmt.Sprintf("sections/%s/%s_citations.json", unqn, unqn)}}

		if val, ok := ignoredNodes[node.Name]; ok && val {
			continue
		}

		if unsafe, err = section.Html(); err != nil {
			return []common.Section{}, []common.Citations{}, err
		}

		node.ID = fmt.Sprintf("sections/%s/%s", unqn, unqn)
		ct.IsPartOf = []common.Entity{common.Entity{Identifier: node.ID}}
		node.Text = unsafe
		nodes = append(nodes, node)
		cits = append(cits, ct)
	}

	return nodes, cits, nil
}
