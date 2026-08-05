package response

type PageResponse struct {
	PageTitle string `json:"page_title"`
	Data      any    `json:"data"`
}

type ErrorPageResponse struct {
	PageTitle string `json:"page_title"`
	Message   string `json:"message"`
}
