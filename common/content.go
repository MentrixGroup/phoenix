package common

import (
	"time"
)

//Language represents human language
type Language struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Direction  string `json:"direction"`
}

//Project represents single language wiki
type Project struct {
	Name       string   `json:"name"`
	Identifier string   `json:"identifier"`
	InLanguage Language `json:"in_language"`
	URL        string   `json:"url"`
	Size       struct {
		Value    int    `json:"value"`
		UnitText string `json:"unit_text"`
	} `json:"size"`
}

//License represents a content Licence
type License struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	URL        string `json:"url"`
}

//Entity represents a subject of a page
type Entity struct {
	Identifier string `json:"identifier"`
}

//Section represents a section of the page

type Section struct {
	Name           string    `json:"name"`
	Identifier     string    `json:"identifier"`
	Version        int       `json:"version"`
	IsPartOf       []Entity  `json:"is_part_of"`
	Text           string    `json:"text"`
	EncodingFormat string    `json:"encoding_format"`
	License        []License `json:"license"`
	Citation       []string  `json:"citation"`
}

// Page represents the root node of a document graph
type Page struct {
	Name         string            `json:"name"`
	Identifier   int               `json:"identifier"`
	URL          string            `json:"url"`
	InLanguage   Language          `json:"in_language"`
	IsPartOf     []Project         `json:"is_part_of"`
	Version      int               `json:"version"`
	DateModified time.Time         `json:"date_modified"`
	License      []License         `json:"license"`
	MainEntity   Entity            `json:"main_entity"`
	About        map[string]string `json:"about"`
	Keywords     string            `json:"keywords"`
	HasPart      []Entity          `json:"has_part"`
}

// Source represents information on the source of the document.
type Source struct {
	// Page ID according to MediaWiki
	ID int `json:"id"`

	// Revision ID according to MediaWiki
	Revision int `json:"revision"`

	// Type 1 UUID; Date and time the source document was rendered
	TimeUUID string `json:"tid"`

	// The wiki/project/hostname of source document
	Authority string `json:"authority"`
}

type metadata struct {
	ID      string `json:"-"`
	Context string `json:"@context"`
	Type    string `json:"@type"`
}

// Thing corresponds to https://schema.org/Thing
type Thing struct {
	metadata
	AlternateName string `json:"alternateName,omitempty"`
	Description   string `json:"description,omitempty"`
	Image         string `json:"image,omitempty"`
	Name          string `json:"name,omitempty"`
	SameAs        string `json:"sameAs"`
}

// NewThing returns an initialized Thing
func NewThing() *Thing {
	return &Thing{metadata: metadata{Context: "https://schema.org", Type: "Thing"}}
}

//GetEnLang returns en Language model
func GetEnLang() *Language {
	return &Language{
		Name:       "English",
		Identifier: "en",
		Direction:  "ltr",
	}
}

func NewLicense() *License {
	return &License{
		Identifier: "CC-BY-SA-3.0",
		Name:       "Creative Commons Attribution Share Alike 3.0 Unported",
		URL:        "https://creativecommons.org/licenses/by-sa/3.0/",
	}
}

// RelatedTopic corresponds to a Wikidata topic that relates to a Node's content.
type RelatedTopic struct {
	ID       string  `json:"id"`
	Label    string  `json:"label"`
	Salience float32 `json:"salience"`
}

// Citation single citation data object
type Citation struct {
	Identifier string   `json:"identifier"`
	Text       string   `json:"text"`
	References []string `json:"references"`
}

// Citations citations collection
type Citations struct {
	Citations []Citation `json:"citations"`
	IsPartOf  []Entity   `json:"is_part_of"`
}

// Book book source
type Book struct {
	Isbn          string
	Name          string
	Author        []string
	Publisher     string
	Datepublished string
	Thumbnailurl  string
}
