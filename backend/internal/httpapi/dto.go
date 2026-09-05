package httpapi

type ConvertRequest struct {
	Pseudocode string `json:"pseudocode"`
}

type ConvertResponse struct {
	GoCode string `json:"go_code,omitempty"`
	Error string `json:"error,omitempty"`
}