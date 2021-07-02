package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/wikimedia/phoenix/common"
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

func getCitation(document *goquery.Document, cId, page string) (common.Citation, error) {
	var err error
	var unsafe string
	citation := common.Citation{}

	slct := document.Find(fmt.Sprintf("#%s", cId))

	if len(slct.Nodes) <= 0 {
		return citation, fmt.Errorf(fmt.Sprintf("No element was found with id=%s found", cId))
	}
	li := slct.Eq(0)

	if unsafe, err = li.Html(); err != nil {
		return citation, fmt.Errorf("error during getting list item HTML content")
	}

	citation.Identifier = li.AttrOr("id", "")
	citation.References = getRefs(li, page)
	citation.Text = unsafe

	return citation, nil
}

func parseParsoidDocumentNodes(document *goquery.Document, page *common.Page) ([]common.Node, []common.Citations, error) {
	var err error
	var modified = page.DateModified
	var nameCounts = make(map[string]int)
	var nodes = make([]common.Node, 0)
	var sections = document.Find("html>body>section[data-mw-section-id]")
	var cits = make([]common.Citations, 0)

	for i := range sections.Nodes {
		var node = common.Node{}
		var ct = common.Citations{}
		var section = sections.Eq(i)
		var unsafe string
		var nodeCites = make([]common.Citation, 0)

		for _, id := range getCitationIds(section, page.Name) {
			ct, err := getCitation(document, id, page.Name)

			if err != nil {
				fmt.Println(fmt.Sprintf("problem with creating citation: %s", err))
			}
			nodeCites = append(nodeCites, ct)
		}

		node.Source = page.Source

		ct.Citations = nodeCites
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

		node.DateModified = modified

		node.HasPart = []string{fmt.Sprintf("nodes/%s_%s_citations.json", replaceSpaces(page.Name), replaceSpaces(node.Name))}

		if val, ok := ignoredNodes[node.Name]; ok && val {
			continue
		}

		if unsafe, err = section.Html(); err != nil {
			return []common.Node{}, []common.Citations{}, err
		}

		node.ID = fmt.Sprintf("nodes/%s_%s", replaceSpaces(page.Name), replaceSpaces(node.Name))
		ct.IsPartOf = []string{node.ID}
		node.Unsafe = unsafe
		nodes = append(nodes, node)

		cits = append(cits, ct)
	}

	return nodes, cits, nil
}
