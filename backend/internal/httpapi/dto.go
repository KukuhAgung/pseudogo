package httpapi

type ConvertRequest struct {
	Pseudocode string `json:"pseudocode"`
}

type ConvertResponse struct {
	GoCode string `json:"go_code,omitempty"`
	Error string `json:"error,omitempty"`
}

type RunRequest struct {
	Pseudocode string `json:"pseudocode"`
	Input      string `json:"input"`
}

type RunResponse struct {
	Output   string `json:"output,omitempty"`
	Stderr   string `json:"stderr,omitempty"`
	TimedOut bool   `json:"timed_out,omitempty"`
	Error    string `json:"error,omitempty"`
}