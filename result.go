package gocsfixer

type Result struct {
	Type string `json:"type"`
	File string `json:"file"`
	Line int `json:"line"`
	Text string `json:"text"`
}