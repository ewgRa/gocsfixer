package response

type DiffLinerResponse struct {
	File string `json:"file"`
	Line int `json:"line"`
	DiffLine int `json:"diff_line"`
}