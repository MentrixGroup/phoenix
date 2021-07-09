package common

// ChangeEvent is a JSON object representing a change to a document
type ChangeEvent struct {
	Title      string `json:"title"`
	ServerName string `json:"server_name"`
	Revision   int    `json:"revision"`
}

// NodeStoredEvent is a JSON object that corresponds to a Node being added to the Content Store.
type NodeStoredEvent struct {
	ID string `json:"id"`
}

// SourseParseEvent is a JSON object that corresponds to a Source being found in citation.
type SourseParseEvent struct {
	ID      string `json:"id"`
	Page    string `json:"page"`
	Section string `json:"section"`
}
