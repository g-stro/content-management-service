package response

type CreateContent struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	CreationDate string `json:"created_at"`
}

type GetContent struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Details     []Details `json:"details"`
}

type Details struct {
	ContentType string `json:"content_type"`
	Value       string `json:"value"`
}
