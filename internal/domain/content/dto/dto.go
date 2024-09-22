package dto

type ContentWithDetailsRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Details     []Detail `json:"details"`
}

type ContentWithDetailsResponse struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Details     []Detail `json:"details"`
}

type Detail struct {
	ContentType string `json:"content_type"`
	Value       string `json:"value"`
}
