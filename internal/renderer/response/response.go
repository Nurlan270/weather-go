package response

type ErrorMessages map[string]string
type OldData map[string]string

type PageResponse struct {
	PageTitle string    `json:"page_title"`
	Error     ErrorData `json:"error,omitempty"`
	Data      any       `json:"data,omitempty"`
	OldData   OldData   `json:"old_data,omitempty"`
}

type ErrorData struct {
	Message string        `json:"message"`
	Items   ErrorMessages `json:"items,omitempty"`
}

type ErrorPageResponse struct {
	PageTitle string `json:"page_title"`
	Message   string `json:"message,omitempty"`
}
